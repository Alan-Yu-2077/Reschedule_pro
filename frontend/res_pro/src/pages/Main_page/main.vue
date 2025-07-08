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
                <button class="slot-btn">{{ getCourseName(i, j) }}</button>
              </td>
            </tr>
          </tbody>
        </table>
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