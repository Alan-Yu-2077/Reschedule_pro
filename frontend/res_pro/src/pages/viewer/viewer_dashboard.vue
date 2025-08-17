<template>
  <div class="main-container">
    <!-- 侧边栏切换按钮 -->
    <div class="sidebar-toggle" @click="toggleSidebar" title="Toggle Sidebar">
      <div class="hamburger-icon">
        <span></span>
        <span></span>
        <span></span>
      </div>
      <span class="toggle-text">Menu</span>
    </div>

    <!-- 背景遮罩 -->
    <div class="sidebar-overlay" :class="{ active: sidebarExpanded }" @click="closeSidebar"></div>

    <!-- 侧边栏 -->
    <div class="sidebar" :class="{ 'sidebar-expanded': sidebarExpanded }">
      <div class="class-header">
        <h2>Class List</h2>
        <!-- 只读用户不显示添加按钮 -->
      </div>
      <ul>
        <li v-for="cls in classList" :key="cls" @click="selectClass(cls)">
          {{ cls }}
        </li>
      </ul>
      <!-- 登出按钮放在侧边栏底部 -->
      <div class="logout-section">
        <button class="logout-btn" @click="logout">
          <span class="logout-icon">🚪</span>
          <span class="logout-text">Logout</span>
        </button>
      </div>
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
                <div class="slot-display">
                  {{ getCourseName(i, j) }}
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- 只读用户不显示课程调换面板 -->
      
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
import { onLoad } from '@dcloudio/uni-app';

const classList = ref([]);
const currentClass = ref('');
const logs = ref([]);
const scheduleData = ref([]);
const sidebarExpanded = ref(false); // 侧边栏默认收起

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
  // 选择班级后自动收起侧边栏
  sidebarExpanded.value = false;
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

const logout = () => {
  uni.clearStorageSync(); // 清除本地存储
  uni.showToast({ title: 'Logged out successfully', icon: 'success' });
  setTimeout(() => {
    uni.reLaunch({
      url: '/pages/login/login'
    });
  }, 1000);
};

const toggleSidebar = () => {
  sidebarExpanded.value = !sidebarExpanded.value;
};

const closeSidebar = () => {
  sidebarExpanded.value = false;
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
  position: relative;
  overflow: hidden; /* 防止侧边栏溢出 */
}

.sidebar-toggle {
  position: fixed;
  top: 20px;
  left: 20px;
  z-index: 1000;
  cursor: pointer;
  padding: 12px;
  background: #fff;
  border: 2px solid #333;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  display: flex; /* 始终显示汉堡菜单按钮 */
  align-items: center; /* 垂直居中 */
  gap: 8px; /* 按钮和文字之间的间距 */
  transition: all 0.3s ease;
}

.sidebar-toggle:hover {
  background: #f8f9fa;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
  transform: scale(1.1);
  border-color: #007bff;
}

.hamburger-icon {
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  width: 20px;
  height: 16px;
}

.hamburger-icon span {
  display: block;
  width: 100%;
  height: 2px;
  background: #333;
  border-radius: 2px;
  transition: all 0.3s ease;
}

/* 汉堡菜单动画效果 */
.sidebar-toggle:hover .hamburger-icon span:nth-child(1) {
  transform: rotate(45deg) translate(5px, 5px);
}

.sidebar-toggle:hover .hamburger-icon span:nth-child(2) {
  opacity: 0;
}

.sidebar-toggle:hover .hamburger-icon span:nth-child(3) {
  transform: rotate(-45deg) translate(7px, -6px);
}

.toggle-text {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  letter-spacing: 0.5px;
}

.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: 250px;
  height: 100vh;
  border-right: 2px dashed #000;
  padding: 20px;
  padding-top: 80px; /* 增加顶部内边距，避免被汉堡菜单按钮遮挡 */
  display: flex;
  flex-direction: column;
  background: #fdfdfb;
  z-index: 200;
  transform: translateX(-100%);
  transition: transform 0.3s ease-in-out;
  box-shadow: 2px 0 10px rgba(0, 0, 0, 0.1);
}

.sidebar-expanded {
  transform: translateX(0);
}

.sidebar-expanded {
  transform: translateX(0);
}

.sidebar h2 {
  margin-bottom: 10px;
}

.sidebar ul {
  list-style: none;
  padding: 0;
  flex: 1;
  overflow-y: auto;
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

.logout-section {
  margin-top: auto;
  padding-top: 20px;
  border-top: 1px solid #e0e0e0;
}

.logout-btn {
  width: 100%;
  padding: 12px 16px;
  background: linear-gradient(135deg, #ff6b6b, #ee5a52);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(255, 107, 107, 0.3);
}

.logout-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 107, 107, 0.4);
  background: linear-gradient(135deg, #ff5252, #d32f2f);
}

.logout-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 4px rgba(255, 107, 107, 0.3);
}

.logout-icon {
  font-size: 16px;
}

.logout-text {
  font-size: 14px;
  letter-spacing: 0.5px;
}

.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 4;
  opacity: 0;
  visibility: hidden;
  transition: opacity 0.3s ease-in-out, visibility 0.3s ease-in-out;
}

.sidebar-overlay.active {
  opacity: 1;
  visibility: visible;
}

.schedule-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 20px;
  min-width: 0;
  width: 100%;
  height: 100vh;
  box-sizing: border-box;
  transition: all 0.3s ease;
  max-width: 100vw;
  padding: 0;
}

/* 响应式设计：在小屏幕上调整布局 */
@media (max-width: 768px) {
  .main-container {
    padding: 4px;
    height: 100vh;
    overflow-y: auto;
    overflow-x: hidden;
  }
  
  .sidebar {
    width: 280px;
  }
  
  .schedule-panel {
    padding: 4px;
    height: auto;
    min-height: 100vh;
    overflow: visible;
    display: flex;
    flex-direction: column;
  }
  
  .schedule-top {
    flex-direction: column;
    gap: 4px;
    margin-bottom: 8px;
    flex-shrink: 0;
  }
  
  .schedule-top .title h2 {
    font-size: 1em;
    margin: 0;
    line-height: 1.2;
  }
  
  .schedule-body {
    flex: 1;
    min-height: 0;
    padding: 4px 0;
    overflow: visible;
    border: 1px dashed #000;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .schedule-table {
    width: 90%;
    max-width: 400px;
    font-size: 0.8em;
    table-layout: fixed;
  }
  
  .schedule-table th,
  .schedule-table td {
    width: auto;
    min-width: 20px;
    font-size: 0.8em;
    padding: 2px;
  }
  
  .schedule-table th {
    height: 30px;
    min-height: 30px;
    padding: 2px;
  }
  
  .schedule-table td {
    height: 30px;
    min-height: 30px;
  }
  
  .schedule-table td:first-child {
    width: 50px;
    min-width: 50px;
    padding: 2px;
    font-size: 0.7em;
  }
  
  .slot-display {
    width: 25px !important;
    height: 30px !important;
    min-width: 25px !important;
    min-height: 30px !important;
    padding: 2px;
    font-size: 0.7em;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .log-panel {
    max-height: 80px;
    padding: 5px;
    flex-shrink: 0;
    overflow: auto;
  }
  
  .log-panel h3 {
    font-size: 0.9em;
    margin-bottom: 5px;
  }
  
  .log-content {
    font-size: 0.7em;
  }
}

/* 超小屏幕优化 */
@media (max-width: 480px) {
  .main-container {
    padding: 2px;
  }
  
  .schedule-panel {
    padding: 2px;
  }
  
  .schedule-top {
    gap: 2px;
    margin-bottom: 4px;
  }
  
  .schedule-top .title h2 {
    font-size: 0.9em;
  }
  
  .schedule-body {
    padding: 2px 0;
  }
  
  .schedule-table {
    width: 100%;
    font-size: 0.7em;
  }
  
  .schedule-table th,
  .schedule-table td {
    width: auto;
    min-width: 20px;
    font-size: 0.7em;
    padding: 2px;
  }
  
  .schedule-table th {
    height: 25px;
    min-height: 25px;
    padding: 25px;
  }
  
  .schedule-table td {
    height: 25px;
    min-height: 25px;
  }
  
  .schedule-table td:first-child {
    width: 40px;
    min-width: 40px;
    padding: 2px;
    font-size: 0.6em;
  }
  
  .slot-display {
    width: 20px !important;
    height: 25px !important;
    min-width: 20px !important;
    min-height: 25px !important;
    padding: 2px;
    font-size: 0.6em;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .log-panel {
    max-height: 60px;
    padding: 3px;
  }
  
  .log-panel h3 {
    font-size: 0.8em;
    margin-bottom: 3px;
  }
  
  .log-content {
    font-size: 0.6em;
  }
}

.class-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
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

.schedule-table {
  width: 100%;
  max-width: 100vw;
  margin: 0 auto;
  border-collapse: collapse;
  font-size: 1em;
  table-layout: fixed;
  text-align: center;
}

.schedule-table th,
.schedule-table td {
  border: 1px solid #000;
  padding: 0;
  font-size: 0.8em;
  width: 40px;
  min-width: 40px;
  height: 50px;
  min-height: 50px;
}

.schedule-table thead {
  background-color: #f5f5f5;
}

.schedule-table th {
  background-color: #eaeaea;
  font-weight: bold;
  font-size: 0.8em;
  height: 50px;
  min-height: 50px;
  padding: 8px;
}

.schedule-table td {
  height: 50px;
  min-height: 50px;
  vertical-align: middle;
  padding: 0;
}

.schedule-table td:first-child {
  font-weight: bold;
  background-color: #f8f8f8;
  width: 60px;
  min-width: 60px;
  padding: 0;
  font-size: 0.95em;
}

.slot-display {
  width: 40px !important;
  min-width: 40px !important;
  height: 50px !important;
  min-height: 50px !important;
  font-size: 0.8em;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  -webkit-tap-highlight-color: transparent; /* 移除移动端点击高亮 */
  touch-action: manipulation; /* 优化触摸操作 */
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
</style> 