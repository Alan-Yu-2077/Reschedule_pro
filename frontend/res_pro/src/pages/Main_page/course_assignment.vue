<template>
  <div class="main-container">
    <div class="title-area">
      <h2>Course Assignment</h2>
      <button class="back-btn" @click="goBack">← Back</button>
    </div>

    <!-- 课程信息输入区域 -->
    <div class="course-input-section">
      <div class="input-group">
        <label>Course Name:</label>
        <textarea 
          v-model="courseName" 
          placeholder="Enter course name" 
          class="course-input"
          auto-height
          maxlength="30"
        />
      </div>

      <!-- 周数类型选择 -->
      <div class="week-type-selection">
        <label>Week Type:</label>
        <div class="week-type-buttons">
          <button 
            class="week-type-btn" 
            :class="{ active: weekType === 'continuous' }"
            @click="weekType = 'continuous'"
          >
            Continuous Weeks
          </button>
          <button 
            class="week-type-btn" 
            :class="{ active: weekType === 'discrete' }"
            @click="weekType = 'discrete'"
          >
            Discrete Weeks
          </button>
        </div>
      </div>

      <!-- 连续周选择 -->
      <div v-if="weekType === 'continuous'" class="continuous-weeks">
        <div class="week-range">
          <div class="range-input">
            <label>Start Week:</label>
            <select v-model="startWeek" class="week-select">
              <option v-for="week in 20" :key="week" :value="week">{{ week }}</option>
            </select>
          </div>
          <div class="range-input">
            <label>End Week:</label>
            <select v-model="endWeek" class="week-select">
              <option v-for="week in 20" :key="week" :value="week">{{ week }}</option>
            </select>
          </div>
        </div>
      </div>

      <!-- 离散周选择 -->
      <div v-if="weekType === 'discrete'" class="discrete-weeks">
        <label>Select Weeks:</label>
        <div class="week-grid">
          <button 
            v-for="week in 20" 
            :key="week"
            class="week-number-btn"
            :class="{ selected: selectedWeeks.includes(week) }"
            @click="toggleWeek(week)"
          >
            {{ week }}
          </button>
        </div>
      </div>
    </div>

    <!-- 课程表分配区域 -->
    <div class="schedule-section">
      <h3>Click on time slots to assign the course:</h3>
      <div class="schedule-body">
        <table class="schedule-table">
          <thead>
            <tr>
              <th>Period</th>
              <th v-for="day in ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']" :key="day">{{ day }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(period, i) in ['Morning Slot 1', 'Morning Slot 2', 'Afternoon Slot 1', 'Afternoon Slot 2', 'Evening Slot']" :key="i">
              <td>{{ period }}</td>
              <td v-for="(day, j) in 7" :key="j">
                <button 
                  class="slot-btn" 
                  :class="{ assigned: isAssigned(i, j) }"
                  @click="assignCourse(i, j)"
                >
                  {{ getSlotDisplay(i, j) }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 底部按钮区域 -->
    <div class="bottom-controls">
      <button class="control-btn" @click="clearAssignments">Clear All</button>
      <button class="control-btn" @click="saveAssignments">Save & Return</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'

const courseName = ref('')
const weekType = ref('continuous')
const startWeek = ref(1)
const endWeek = ref(1)
const selectedWeeks = ref([])

// 存储分配的课程
const assignments = ref([])

onShow(() => {
  courseName.value = ''
  weekType.value = 'continuous'
  startWeek.value = 1
  endWeek.value = 1
  selectedWeeks.value = []
  assignments.value = []
})

const goBack = () => {
  uni.navigateBack()
}

const toggleWeek = (week) => {
  const index = selectedWeeks.value.indexOf(week)
  if (index > -1) {
    selectedWeeks.value.splice(index, 1)
  } else {
    selectedWeeks.value.push(week)
  }
}

const isAssigned = (row, col) => {
  return assignments.value.some(assignment => 
    assignment.row === row && assignment.col === col
  )
}

const getSlotDisplay = (row, col) => {
  const assignment = assignments.value.find(a => a.row === row && a.col === col)
  return assignment ? assignment.courseName : ''
}

const assignCourse = (row, col) => {
  if (!courseName.value.trim()) {
    uni.showToast({
      title: 'Please enter a course name',
      icon: 'none',
      duration: 2000
    })
    return
  }

  if (weekType.value === 'continuous') {
    if (startWeek.value > endWeek.value) {
      uni.showToast({
        title: 'Start week cannot be greater than end week',
        icon: 'none',
        duration: 2000
      })
      return
    }
  } else if (weekType.value === 'discrete') {
    if (selectedWeeks.value.length === 0) {
      uni.showToast({
        title: 'Please select at least one week',
        icon: 'none',
        duration: 2000
      })
      return
    }
  }

  // 检查是否已经分配
  const existingIndex = assignments.value.findIndex(a => a.row === row && a.col === col)
  
  if (existingIndex > -1) {
    // 移除现有分配
    assignments.value.splice(existingIndex, 1)
    uni.showToast({
      title: 'Course assignment removed',
      icon: 'success',
      duration: 1500
    })
  } else {
    // 添加新分配
    const assignment = {
      row,
      col,
      courseName: courseName.value,
      weekType: weekType.value,
      startWeek: startWeek.value,
      endWeek: endWeek.value,
      selectedWeeks: [...selectedWeeks.value]
    }
    
    assignments.value.push(assignment)
    
    uni.showToast({
      title: `Course "${courseName.value}" assigned successfully!`,
      icon: 'success',
      duration: 2000
    })
  }
}

const clearAssignments = () => {
  assignments.value = []
  uni.showToast({
    title: 'All assignments cleared',
    icon: 'success',
    duration: 1500
  })
}

const saveAssignments = () => {
  if (assignments.value.length === 0) {
    uni.showToast({
      title: 'No assignments to save',
      icon: 'none',
      duration: 2000
    })
    return
  }

  // 这里可以将assignments保存到全局状态或本地存储
  // 暂时使用uni-app的本地存储
  try {
    uni.setStorageSync('courseAssignments', assignments.value)
    uni.showToast({
      title: 'Assignments saved successfully!',
      icon: 'success',
      duration: 2000
    })
    
    // 返回上一页
    setTimeout(() => {
      uni.navigateBack()
    }, 2000)
  } catch (e) {
    uni.showToast({
      title: 'Failed to save assignments',
      icon: 'none',
      duration: 2000
    })
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Patrick+Hand&display=swap');

.main-container {
  font-family: 'Patrick Hand', cursive;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  background: #fdfdfb;
  background-image: url('https://www.transparenttextures.com/patterns/paper-fibers.png');
  background-size: auto;
  min-height: 100vh;
  box-sizing: border-box;
}

.title-area {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  margin-bottom: 30px;
}

.title-area h2 {
  color: #333;
  font-size: 2em;
  margin: 0;
}

.back-btn {
  padding: 10px 20px;
  background: white;
  border: 2px solid #000;
  border-radius: 8px;
  font-family: 'Patrick Hand', cursive;
  font-size: 1.1em;
  cursor: pointer;
  transition: all 0.2s;
}

.back-btn:hover {
  background: #f5f5f5;
  transform: translateY(-2px);
}

.course-input-section {
  width: 100%;
  max-width: 600px;
  margin-bottom: 30px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  border: 2px solid #ddd;
}

.input-group {
  margin-bottom: 20px;
}

.input-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  color: #333;
}

.course-input {
  width: 100%;
  padding: 12px;
  border: 2px solid #ddd;
  border-radius: 8px;
  font-family: 'Patrick Hand', cursive;
  font-size: 1em;
  box-sizing: border-box;
  resize: none;
  min-height: 20px;
  line-height: 1.2;
}

.course-input:focus {
  outline: none;
  border-color: #007bff;
}

.week-type-selection {
  margin-bottom: 20px;
}

.week-type-selection label {
  display: block;
  margin-bottom: 10px;
  font-weight: bold;
  color: #333;
}

.week-type-buttons {
  display: flex;
  gap: 10px;
}

.week-type-btn {
  flex: 1;
  padding: 12px;
  border: 2px solid #ddd;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 1em;
  transition: all 0.2s;
}

.week-type-btn.active {
  border-color: #007bff;
  background: #e3f2fd;
  color: #007bff;
}

.week-type-btn:hover {
  border-color: #007bff;
}

.continuous-weeks {
  margin-bottom: 20px;
}

.week-range {
  display: flex;
  gap: 15px;
}

.range-input {
  flex: 1;
}

.range-input label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  color: #333;
}

.week-select {
  width: 100%;
  padding: 10px;
  border: 2px solid #ddd;
  border-radius: 8px;
  font-family: 'Patrick Hand', cursive;
  font-size: 1em;
  background: white;
}

.discrete-weeks {
  margin-bottom: 20px;
}

.discrete-weeks label {
  display: block;
  margin-bottom: 10px;
  font-weight: bold;
  color: #333;
}

.week-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 8px;
  max-width: 300px;
}

.week-number-btn {
  aspect-ratio: 1;
  width: 100%;
  padding: 0;
  border: 2px solid #ddd;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 1em;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 40px;
}

.week-number-btn.selected {
  border-color: #007bff;
  background: #007bff;
  color: white;
}

.week-number-btn:hover {
  border-color: #007bff;
  background: #e3f2fd;
}

.schedule-section {
  width: 100%;
  margin-bottom: 30px;
}

.schedule-section h3 {
  text-align: center;
  color: #333;
  margin-bottom: 20px;
}

.schedule-body {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  max-width: 800px;
  border: 2px dashed #000;
  padding: 20px;
  background: white;
  border-radius: 12px;
  margin: 0 auto;
}

.schedule-table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
  text-align: center;
}

.schedule-table th,
.schedule-table td {
  border: 1px solid #000;
  padding: 10px 6px;
}

.schedule-table th {
  background-color: #eaeaea;
  font-weight: bold;
  font-size: 0.9em;
  padding: 8px 6px;
}

.schedule-table td:first-child {
  font-weight: bold;
  background-color: #f8f8f8;
  width: 100px;
  font-size: 0.85em;
}

.slot-btn {
  width: 100%;
  height: 100%;
  padding: 6px 3px;
  background-color: #ffffff;
  border: 1px solid #ccc;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 0.8em;
  min-height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  border-radius: 4px;
  word-wrap: break-word;
  line-height: 1.2;
}

.slot-btn:hover {
  background-color: #f0f0f0;
  transform: scale(1.02);
}

.slot-btn.assigned {
  background-color: #e3f2fd;
  border-color: #007bff;
  color: #007bff;
  font-weight: bold;
}

.bottom-controls {
  display: flex;
  gap: 20px;
  justify-content: center;
  width: 100%;
  max-width: 600px;
  margin-top: auto;
}

.control-btn {
  padding: 15px 30px;
  font-family: 'Patrick Hand', cursive;
  font-size: 1.1em;
  cursor: pointer;
  border: 2px solid #000;
  background: #fff;
  transition: all 0.2s;
  border-radius: 8px;
  min-width: 120px;
  font-weight: bold;
}

.control-btn:hover {
  background: #f5f5f5;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-container {
    padding: 15px;
  }
  
  .title-area h2 {
    font-size: 1.5em;
  }
  
  .week-type-buttons {
    flex-direction: column;
  }
  
  .week-range {
    flex-direction: column;
  }
  
  .week-grid {
    grid-template-columns: repeat(4, 1fr);
  }
  
  .bottom-controls {
    flex-direction: column;
    gap: 15px;
  }
  
  .control-btn {
    width: 100%;
  }
}
</style> 