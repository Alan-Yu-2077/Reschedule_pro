<template>
  <div class="main-container">
    <div class="title-area">
      <h2>Register New Class Schedule</h2>
      <div class="class-name-block">
        <label>Class Name:</label>
        <textarea 
          v-model="className" 
          placeholder="e.g. Grade 23 - Class 1" 
          class="class-name-input"
          auto-height
          maxlength="50"
        />
      </div>
    </div>

    <!-- 周数控制区域 -->
    <div class="week-controls">
      <button class="week-btn" @click="prevWeek">←</button>
      <span class="week-display">Week {{ currentWeek }}</span>
      <button class="week-btn" @click="nextWeek">→</button>
    </div>

    <!-- 课程表区域 -->
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
              <button class="slot-btn" @click="handleSlotClick(i, j)">
                {{ getSlotDisplay(i, j) }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 底部按钮区域 -->
    <div class="bottom-controls">
      <button class="control-btn" @click="showAddCourseModal">Add Course</button>
      <button class="control-btn" @click="submitSchedule">Submit</button>
    </div>

    <!-- 待分配课程状态显示 -->
    <div v-if="pendingCourse" class="pending-course-indicator">
      <div class="pending-info">
        <span class="pending-label">Ready to assign:</span>
        <span class="pending-course-name">{{ pendingCourse.name }}</span>
        <span class="pending-weeks">
          ({{ getWeekDisplay(pendingCourse) }})
        </span>
      </div>
      <div class="pending-actions">
        <button class="action-btn done-btn" @click="finishAssignment">Done</button>
        <button class="cancel-assign-btn" @click="cancelPendingCourse">Cancel</button>
      </div>
    </div>

    <!-- Add Course 弹窗 -->
    <div v-if="showModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Add New Course</h3>
          <button class="close-btn" @click="closeModal">×</button>
        </div>
        
        <div class="modal-body">
          <!-- 课程名称输入 -->
          <div class="input-group">
            <label>Course Name:</label>
            <textarea 
              v-model="newCourse.name" 
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
                :class="{ active: newCourse.weekType === 'continuous' }"
                @click="newCourse.weekType = 'continuous'"
              >
                Continuous Weeks
              </button>
              <button 
                class="week-type-btn" 
                :class="{ active: newCourse.weekType === 'discrete' }"
                @click="newCourse.weekType = 'discrete'"
              >
                Discrete Weeks
              </button>
            </div>
          </div>

          <!-- 连续周选择 -->
          <div v-if="newCourse.weekType === 'continuous'" class="continuous-weeks">
            <div class="week-range">
              <div class="range-input">
                <label>Start Week:</label>
                <select v-model="newCourse.startWeek" class="week-select">
                  <option v-for="week in 20" :key="week" :value="week">{{ week }}</option>
                </select>
              </div>
              <div class="range-input">
                <label>End Week:</label>
                <select v-model="newCourse.endWeek" class="week-select">
                  <option v-for="week in 20" :key="week" :value="week">{{ week }}</option>
                </select>
              </div>
            </div>
          </div>

          <!-- 离散周选择 -->
          <div v-if="newCourse.weekType === 'discrete'" class="discrete-weeks">
            <label>Select Weeks:</label>
            <div class="week-grid">
              <button 
                v-for="week in 20" 
                :key="week"
                class="week-number-btn"
                :class="{ selected: newCourse.selectedWeeks.includes(week) }"
                @click="toggleWeek(week)"
              >
                {{ week }}
              </button>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="modal-btn cancel-btn" @click="closeModal">Cancel</button>
          <button class="modal-btn confirm-btn" @click="confirmAddCourse">Add Course</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'

const className = ref('')
const currentWeek = ref(1)
const schedule = ref(
  Array(5).fill(null).map(() => Array(7).fill(null))
)

// 弹窗相关状态
const showModal = ref(false)
const newCourse = ref({
  name: '',
  weekType: 'continuous',
  startWeek: 1,
  endWeek: 1,
  selectedWeeks: []
})

// 临时存储待分配的课程
const pendingCourse = ref(null)

onShow(() => {
  className.value = ''
  schedule.value = Array(5).fill(null).map(() => Array(7).fill(null))
  currentWeek.value = 1
  resetNewCourse()
  pendingCourse.value = null
})

const resetNewCourse = () => {
  newCourse.value = {
    name: '',
    weekType: 'continuous',
    startWeek: 1,
    endWeek: 1,
    selectedWeeks: []
  }
}

const prevWeek = () => {
  if (currentWeek.value > 1) currentWeek.value--;
}

const nextWeek = () => {
  if (currentWeek.value < 20) currentWeek.value++;
}

const getSlotDisplay = (row, col) => {
  const cell = schedule.value[row][col];
  if (cell && cell.name) {
    // 检查当前周是否在课程的周数范围内
    if (cell.weekType === 'continuous') {
      if (currentWeek.value >= cell.startWeek && currentWeek.value <= cell.endWeek) {
        return cell.name;
      }
    } else if (cell.weekType === 'discrete') {
      if (cell.selectedWeeks.includes(currentWeek.value)) {
        return cell.name;
      }
    }
  }
  return '';
}

const handleSlotClick = (row, col) => {
  if (pendingCourse.value) {
    // 分配课程到该时间段
    schedule.value[row][col] = {
      name: pendingCourse.value.name,
      weekType: pendingCourse.value.weekType,
      startWeek: pendingCourse.value.startWeek,
      endWeek: pendingCourse.value.endWeek,
      selectedWeeks: [...pendingCourse.value.selectedWeeks]
    };
    
    // 不清空待分配的课程，允许继续分配
    // pendingCourse.value = null;
    
    // 显示成功消息
    uni.showToast({
      title: `Course "${schedule.value[row][col].name}" assigned to time slot!`,
      icon: 'success',
      duration: 1500
    });
  } else {
    // 如果没有待分配的课程，显示提示
    uni.showToast({
      title: 'Please add a course first',
      icon: 'none',
      duration: 2000
    });
  }
}

const showAddCourseModal = () => {
  resetNewCourse()
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  resetNewCourse()
}

const toggleWeek = (week) => {
  const index = newCourse.value.selectedWeeks.indexOf(week)
  if (index > -1) {
    newCourse.value.selectedWeeks.splice(index, 1)
  } else {
    newCourse.value.selectedWeeks.push(week)
  }
}

const confirmAddCourse = () => {
  if (!newCourse.value.name.trim()) {
    uni.showToast({
      title: 'Please enter a course name',
      icon: 'none',
      duration: 2000
    });
    return
  }

  if (newCourse.value.weekType === 'continuous') {
    if (newCourse.value.startWeek > newCourse.value.endWeek) {
      uni.showToast({
        title: 'Start week cannot be greater than end week',
        icon: 'none',
        duration: 2000
      });
      return
    }
  } else if (newCourse.value.weekType === 'discrete') {
    if (newCourse.value.selectedWeeks.length === 0) {
      uni.showToast({
        title: 'Please select at least one week',
        icon: 'none',
        duration: 2000
      });
      return
    }
  }

  // 将课程信息存储到待分配状态
  pendingCourse.value = {
    name: newCourse.value.name,
    weekType: newCourse.value.weekType,
    startWeek: newCourse.value.startWeek,
    endWeek: newCourse.value.endWeek,
    selectedWeeks: [...newCourse.value.selectedWeeks]
  }
  
  closeModal()
  
  // 显示提示信息
  uni.showToast({
    title: `Course "${pendingCourse.value.name}" ready to assign. Click on a time slot to assign it.`,
    icon: 'none',
    duration: 3000
  });
}

const getWeekDisplay = (course) => {
  if (course.weekType === 'continuous') {
    return `${course.startWeek} - ${course.endWeek}`;
  } else if (course.weekType === 'discrete') {
    return course.selectedWeeks.join(', ');
  }
  return '';
}

const cancelPendingCourse = () => {
  pendingCourse.value = null;
  uni.showToast({
    title: 'Pending course assignment cancelled',
    icon: 'none',
    duration: 2000
  });
}

const finishAssignment = () => {
  pendingCourse.value = null;
  uni.showToast({
    title: 'Course assignment completed!',
    icon: 'success',
    duration: 2000
  });
}

const submitSchedule = async () => {
  if (!className.value.trim()) {
    uni.showToast({
      title: 'Please enter a class name',
      icon: 'none',
      duration: 2000
    });
    return;
  }

  // 检查是否有课程分配
  let hasCourses = false;
  for (let row = 0; row < schedule.value.length; row++) {
    for (let col = 0; col < schedule.value[row].length; col++) {
      if (schedule.value[row][col] !== null) {
        hasCourses = true;
        break;
      }
    }
    if (hasCourses) break;
  }

  if (!hasCourses) {
    uni.showToast({
      title: 'Please assign at least one course',
      icon: 'none',
      duration: 2000
    });
    return;
  }

  // 准备提交数据
  const submitData = {
    className: className.value.trim(),
    schedule: schedule.value
  };

  console.log('Submitting data:', submitData);

  try {
    const response = await uni.request({
      url: 'http://localhost:8080/api/schedule/save',
      method: 'POST',
      data: submitData,
      header: {
        'Content-Type': 'application/json'
      }
    });

    console.log('Response:', response);

    if (response.statusCode === 200) {
      uni.showToast({
        title: 'Schedule saved successfully!',
        icon: 'success',
        duration: 2000
      });
      
      // 清空数据，准备下一个班级
      className.value = '';
      schedule.value = Array(5).fill(null).map(() => Array(7).fill(null));
      currentWeek.value = 1;
    } else {
      console.error('Response error:', response);
      throw new Error(`HTTP ${response.statusCode}: ${response.data?.error || 'Unknown error'}`);
    }
  } catch (error) {
    console.error('Submit error:', error);
    uni.showToast({
      title: `Failed to save schedule: ${error.message}`,
      icon: 'none',
      duration: 3000
    });
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
  justify-content: space-between;
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
  flex-direction: column;
  align-items: center;
  margin-bottom: 30px;
  width: 100%;
  max-width: 600px;
}

.title-area h2 {
  margin-bottom: 20px;
  color: #333;
  font-size: 2em;
  text-align: center;
}

.class-name-block {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}

.class-name-block label {
  font-weight: bold;
  color: #333;
  margin-bottom: 8px;
  font-size: 1.1em;
}

.class-name-input {
  padding: 12px 16px;
  width: 100%;
  max-width: 400px;
  font-family: 'Patrick Hand', cursive;
  border: 2px solid #ddd;
  border-radius: 8px;
  font-size: 1em;
  box-sizing: border-box;
  transition: border-color 0.2s;
  resize: none;
  min-height: 20px;
  line-height: 1.2;
}

.class-name-input:focus {
  outline: none;
  border-color: #007bff;
}

.week-controls {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 30px;
  margin-bottom: 30px;
  width: 100%;
}

.week-btn {
  width: 50px;
  height: 50px;
  background: white;
  border: 2px solid #000;
  border-radius: 8px;
  font-size: 1.5em;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  font-family: 'Patrick Hand', cursive;
}

.week-btn:hover {
  background: #f5f5f5;
  transform: scale(1.05);
}

.week-display {
  font-size: 1.4em;
  font-weight: bold;
  color: #333;
  min-width: 120px;
  text-align: center;
}

.schedule-body {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  max-width: 800px;
  border: 2px dashed #000;
  padding: 20px;
  background: white;
  border-radius: 12px;
  margin-bottom: 30px;
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

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.5em;
}

.close-btn {
  background: none;
  border: none;
  font-size: 2em;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
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

.modal-footer {
  display: flex;
  gap: 15px;
  padding: 20px;
  border-top: 1px solid #eee;
  justify-content: flex-end;
}

.modal-btn {
  padding: 12px 24px;
  border: 2px solid #ddd;
  border-radius: 8px;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 1em;
  transition: all 0.2s;
}

.cancel-btn {
  background: white;
  color: #666;
}

.cancel-btn:hover {
  background: #f5f5f5;
  border-color: #999;
}

.confirm-btn {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

.confirm-btn:hover {
  background: #0056b3;
  border-color: #0056b3;
}

/* 待分配课程状态显示样式 */
.pending-course-indicator {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 8px;
  margin-top: 20px;
  background: white;
}

.pending-info {
  display: flex;
  align-items: center;
}

.pending-label {
  font-weight: bold;
  color: #333;
  margin-right: 10px;
}

.pending-course-name {
  color: #333;
}

.pending-weeks {
  color: #666;
}

.pending-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.action-btn {
  padding: 8px 16px;
  border: 2px solid #ddd;
  border-radius: 6px;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 0.9em;
  transition: all 0.2s;
}

.done-btn {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

.done-btn:hover {
  background: #0056b3;
  border-color: #0056b3;
  transform: translateY(-1px);
}

.cancel-assign-btn {
  background: none;
  border: none;
  font-size: 1em;
  cursor: pointer;
  color: #666;
  padding: 0;
  transition: color 0.2s;
}

.cancel-assign-btn:hover {
  color: #333;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-container {
    padding: 15px;
  }
  
  .title-area h2 {
    font-size: 1.5em;
  }
  
  .week-controls {
    gap: 20px;
  }
  
  .week-btn {
    width: 40px;
    height: 40px;
    font-size: 1.2em;
  }
  
  .schedule-body {
    padding: 20px;
  }
  
  .bottom-controls {
    flex-direction: column;
    gap: 15px;
  }
  
  .control-btn {
    width: 100%;
  }

  .modal-content {
    width: 95%;
    margin: 10px;
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
}
</style>