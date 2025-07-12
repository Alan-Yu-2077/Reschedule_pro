<template>
  <div class="main-container">
    <div class="sidebar">
      <div class="class-header">
        <h2>Class List</h2>
        <button class="add-btn" @click="goImportPage">➕</button>
      </div>
      <ul>
        <li v-for="cls in classList" :key="cls" @click="selectClass(cls)">
          {{ cls }}
        </li>
      </ul>
    </div>

    <div class="schedule-panel">
      <div class="schedule-top">
        <div class="controls">
          <button @click="prevWeek">←</button>
        </div>
        <div class="title">
          <h2>{{ currentClass }} - Week {{ currentWeek }} Schedule</h2>
        </div>
        <div class="controls">
          <button @click="nextWeek">→</button>
        </div>
      </div>
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
                  :class="{ 
                    'selected': selectedCourse && selectedCourse.row === i && selectedCourse.col === j,
                    'target': isSelectingTarget && selectedCourse && selectedCourse.row === i && selectedCourse.col === j
                  }"
                  @click="handleSlotClick(i, j)"
                >
                  {{ getCourseName(i, j) }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="course-swap-panel" v-if="selectedCourse">
        <h3>Course Swap</h3>
        <div class="swap-info">
          <p>Selected: {{ getCourseName(selectedCourse.row, selectedCourse.col) }} (Week {{ selectedCourse.sourceWeek }}, {{ getDayName(selectedCourse.col) }}, {{ getPeriodName(selectedCourse.row) }})</p>
          <p>Current Week: {{ currentWeek }} - Click on a slot below to move the course here</p>
          <div class="swap-actions">
            <button @click="cancelSelection" class="cancel-btn">Cancel</button>
            <button @click="deleteCourse" class="delete-btn">Delete Course</button>
            <button @click="startTargetSelection" class="swap-btn">Move to Current Week</button>
          </div>
        </div>
      </div>
      
      <div class="target-selection-panel" v-if="isSelectingTarget">
        <h3>Select Target Position</h3>
        <div class="target-info">
          <p>Click on a slot to move "{{ getCourseName(selectedCourse.row, selectedCourse.col) }}" there</p>
          <button @click="cancelTargetSelection" class="cancel-btn">Cancel</button>
        </div>
      </div>

      <div class="log-panel">
        <h3>Activity Log</h3>
        <div class="log-content">
          <div v-for="(log, index) in logs" :key="index">{{ log }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { onLoad } from '@dcloudio/uni-app';

const classList = ref([]);
const currentClass = ref('');
const logs = ref([]);
const scheduleData = ref([]);
const selectedCourse = ref(null); // 选中的课程
const isSelectingTarget = ref(false); // 是否正在选择目标位置

// 加载班级列表
const loadClasses = async () => {
  try {
    const response = await uni.request({
      url: 'http://localhost:8080/api/schedule/classes',
      method: 'GET',
      header: {
        'Content-Type': 'application/json'
      }
    });

    console.log('Classes API response:', response);

    if (response.statusCode === 200) {
      classList.value = response.data.classes.map(cls => cls.name);
      console.log('Classes loaded:', classList.value);
      
      // 如果有班级，选择第一个
      if (classList.value.length > 0 && !currentClass.value) {
        currentClass.value = classList.value[0];
        loadSchedule();
      }
    } else {
      console.error('Failed to load classes, status:', response.statusCode);
      uni.showToast({ title: 'Failed to load classes', icon: 'none' });
    }
  } catch (error) {
    console.error('Failed to load classes:', error);
    uni.showToast({ title: 'Failed to load classes', icon: 'none' });
  }
};

const selectClass = (cls) => {
  currentClass.value = cls;
  logs.value.push(`Switched to class: ${cls}`);
  loadSchedule();
  loadLogs();
};

const router = useRouter();
const goImportPage = () => {
  uni.navigateTo({
    url: '/pages/Main_page/register_class'
  });
};

const currentWeek = ref(1);
const prevWeek = () => {
  if (currentWeek.value > 1) currentWeek.value--;
  loadSchedule(); // 重新加载当前周的课表
};
const nextWeek = () => {
  if (currentWeek.value < 20) currentWeek.value++;
  loadSchedule(); // 重新加载当前周的课表
};

// 加载课表数据
const loadSchedule = async () => {
  if (!currentClass.value) return;
  
  try {
    const response = await uni.request({
      url: `http://localhost:8080/api/schedule/class/${encodeURIComponent(currentClass.value)}/week/${currentWeek.value}`,
      method: 'GET'
    });

    console.log('Schedule API response:', response);

    if (response.statusCode === 200) {
      scheduleData.value = response.data.schedules || [];
      console.log('Schedule loaded for week', currentWeek.value, ':', scheduleData.value);
    } else {
      console.error('Failed to load schedule, status:', response.statusCode);
      uni.showToast({ title: 'Failed to load schedule', icon: 'none' });
    }
  } catch (error) {
    console.error('Failed to load schedule:', error);
    uni.showToast({ title: 'Failed to load schedule', icon: 'none' });
  }
};

const loadLogs = () => {
  uni.request({
    url: 'http://localhost:8080/admin/logs',
    method: 'GET',
    success: (res) => {
      if (res.statusCode === 200) {
        logs.value = res.data.logs || [];
      }
    },
    fail: () => {
      uni.showToast({ title: 'Failed to load logs', icon: 'none' });
    }
  });
};

// 根据时间段位置获取课程名称
const getCourseName = (row, col) => {
  const match = scheduleData.value.find(schedule => 
    schedule.timeSlotRow === row && 
    schedule.timeSlotCol === col
  );
  
  return match ? match.course.name : '';
};

// 获取星期名称
const getDayName = (col) => {
  const days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
  return days[col] || 'Unknown';
};

// 获取时间段名称
const getPeriodName = (row) => {
  const periods = ['Morning Slot 1', 'Morning Slot 2', 'Afternoon Slot 1', 'Afternoon Slot 2', 'Evening Slot'];
  return periods[row] || 'Unknown';
};

// 处理课程槽点击
const handleSlotClick = (row, col) => {
  if (isSelectingTarget.value) {
    // 正在选择目标位置
    moveCourseToTarget(row, col);
  } else {
    // 选择源课程
    const courseName = getCourseName(row, col);
    if (courseName) {
      selectedCourse.value = { 
        row, 
        col, 
        courseName,
        sourceWeek: currentWeek.value // 保存源周数
      };
    } else {
      uni.showToast({ title: 'No course in this slot', icon: 'none' });
    }
  }
};

// 取消选择
const cancelSelection = () => {
  selectedCourse.value = null;
  isSelectingTarget.value = false;
};

// 删除课程
const deleteCourse = async () => {
  if (!selectedCourse.value) return;
  
  try {
    const response = await uni.request({
      url: 'http://localhost:8080/api/schedule/delete',
      method: 'DELETE',
      header: {
        'Content-Type': 'application/json'
      },
      data: {
        className: currentClass.value,
        weekNumber: selectedCourse.value.sourceWeek, // 删除源周数的课程
        timeSlotRow: selectedCourse.value.row,
        timeSlotCol: selectedCourse.value.col
      }
    });

    if (response.statusCode === 200) {
      uni.showToast({ title: 'Course deleted successfully', icon: 'success' });
      selectedCourse.value = null;
      loadSchedule(); // 重新加载课程表
      loadLogs(); // 重新加载日志
    } else {
      uni.showToast({ title: response.data.error || 'Failed to delete course', icon: 'none' });
    }
  } catch (error) {
    console.error('Failed to delete course:', error);
    uni.showToast({ title: 'Failed to delete course', icon: 'none' });
  }
};

// 开始选择目标位置
const startTargetSelection = () => {
  isSelectingTarget.value = true;
};

// 取消目标选择
const cancelTargetSelection = () => {
  isSelectingTarget.value = false;
};

// 移动课程到目标位置
const moveCourseToTarget = async (targetRow, targetCol) => {
  if (!selectedCourse.value) return;
  
  // 检查目标位置是否已有课程
  const targetCourseName = getCourseName(targetRow, targetCol);
  if (targetCourseName) {
    uni.showToast({ title: 'Target slot is not empty', icon: 'none' });
    return;
  }
  
  // 直接使用当前周数作为目标周数
  const targetWeek = currentWeek.value;
  
  try {
    const response = await uni.request({
      url: 'http://localhost:8080/api/schedule/move',
      method: 'POST',
      header: {
        'Content-Type': 'application/json'
      },
      data: {
        className: currentClass.value,
        sourceWeek: selectedCourse.value.sourceWeek, // 使用选中的源周数
        sourceRow: selectedCourse.value.row,
        sourceCol: selectedCourse.value.col,
        targetWeek: targetWeek,
        targetRow: targetRow,
        targetCol: targetCol
      }
    });

    if (response.statusCode === 200) {
      uni.showToast({ title: 'Course moved successfully', icon: 'success' });
      selectedCourse.value = null;
      isSelectingTarget.value = false;
      loadSchedule(); // 重新加载课程表
      loadLogs(); // 重新加载日志
    } else {
      uni.showToast({ title: response.data.error || 'Failed to move course', icon: 'none' });
    }
  } catch (error) {
    console.error('Failed to move course:', error);
    uni.showToast({ title: 'Failed to move course', icon: 'none' });
  }
};

onMounted(() => {
  loadClasses();
  loadLogs();
});

// 页面显示时刷新数据
onLoad(() => {
  loadClasses();
});
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Patrick+Hand&display=swap');

.main-container {
  display: flex;
  height: 100vh;
  background: #fdfdfb;
  font-family: 'Patrick Hand', cursive;
  padding: 16px;
  box-sizing: border-box;
  background-image: url('https://www.transparenttextures.com/patterns/paper-fibers.png');
  background-size: auto;
}

.sidebar {
  width: 180px;
  border-right: 2px dashed #000;
  padding: 20px;
  flex-shrink: 0;
}

.sidebar h2 {
  margin-bottom: 10px;
}

.sidebar ul {
  list-style: none;
  padding: 0;
}

.sidebar li {
  cursor: pointer;
  padding: 5px;
  border: 1px solid #000;
  margin-bottom: 10px;
  transition: background 0.2s;
}

.sidebar li:hover {
  background: rgba(0, 0, 0, 0.05);
}

.schedule-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 20px;
  min-width: 0;
}

.class-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.add-btn {
  background: none;
  border: none;
  font-size: 1.2em;
  width: 24px;
  height: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s;
}
.add-btn:hover {
  transform: scale(1.2);
}

.schedule-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.schedule-top .title {
  flex: 1;
  text-align: center;
}

.schedule-top .controls button {
  background: none;
  border: none;
  font-size: 1.2em;
  cursor: pointer;
  padding: 2px 6px;
  transition: transform 0.2s;
  line-height: 1;
  height: auto;
}

.schedule-top .controls button:hover {
  transform: scale(1.2);
}

.schedule-body {
  flex: 1;
  border: 2px dashed #000;
  margin-bottom: 10px;
  padding: 15px 0;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.course-swap-panel {
  border-top: 2px dashed #000;
  padding: 10px;
  background: #f8f9fa;
  margin-bottom: 10px;
}

.course-swap-panel h3 {
  margin-bottom: 10px;
  color: #333;
}

.swap-info p {
  margin-bottom: 10px;
  font-weight: bold;
}

.swap-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.swap-actions button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.cancel-btn {
  background: #6c757d;
  color: white;
}

.cancel-btn:hover {
  background: #5a6268;
}

.delete-btn {
  background: #dc3545;
  color: white;
}

.delete-btn:hover {
  background: #c82333;
}

.swap-btn {
  background: #007bff;
  color: white;
}

.swap-btn:hover {
  background: #0056b3;
}

.target-selection-panel {
  border-top: 2px dashed #000;
  padding: 10px;
  background: #fff3cd;
  margin-bottom: 10px;
}

.target-selection-panel h3 {
  margin-bottom: 10px;
  color: #856404;
}

.target-info p {
  margin-bottom: 10px;
  color: #856404;
}

.log-panel {
  border-top: 2px dashed #000;
  padding: 10px;
  max-height: 150px;
  overflow-y: auto;
}

.log-panel h3 {
  margin-bottom: 10px;
}

.log-content {
  font-size: 0.9em;
}

.schedule-table {
  width: 80%;
  max-width: 700px;
  margin: 0 auto;
  border-collapse: collapse;
  font-size: 0.95em;
  table-layout: fixed;
  text-align: center;
}

.schedule-table th,
.schedule-table td {
  border: 1px solid #000;
  padding: 0;
  font-size: 0.95em;
  width: 75px;
  min-width: 75px;
}

.schedule-table thead {
  background-color: #f5f5f5;
}

.schedule-table th {
  background-color: #eaeaea;
  font-weight: bold;
  font-size: 0.95em;
  height: 45px;
  min-height: 45px;
  padding: 8px;
}

.schedule-table td {
  height: 45px;
  min-height: 45px;
  vertical-align: middle;
  padding: 0;
}

.schedule-table td:first-child {
  font-weight: bold;
  background-color: #f8f8f8;
  width: 60px;
  min-width: 60px;
  padding: 8px;
}

.slot-btn {
  width: 100%;
  height: 45px;
  min-height: 45px;
  padding: 6px;
  background-color: #ffffff;
  border: 1px solid #ccc;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  transition: all 0.3s ease;
}

.slot-btn.selected {
  background: #4CAF50;
  color: white;
  border-color: #4CAF50;
}

.slot-btn.target {
  background: #FF9800;
  color: white;
  border-color: #FF9800;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); }
}

.slot-btn {
  font-size: 0.85em;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.slot-btn:hover {
  background-color: #f0f0f0;
}
</style>