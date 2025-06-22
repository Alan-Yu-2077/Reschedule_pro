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

const selectClass = (cls) => {
  currentClass.value = cls;
  logs.value.push(`Switched to class: ${cls}`);
  loadSchedule();
  loadLogs();
};

const router = useRouter();
const goImportPage = () => {
  router.push('/pages/Main_page/register_class');
};

const currentWeek = ref(1);
const prevWeek = () => {
  if (currentWeek.value > 1) currentWeek.value--;
};
const nextWeek = () => {
  currentWeek.value++;
};

const loadSchedule = () => {
  uni.request({
    url: 'http://localhost:8080/schedule/' + currentClass.value,
    method: 'GET',
    success: (res) => {
      if (res.statusCode === 200) {
        scheduleData.value = res.data.courses || [];
        console.log('Schedule loaded:', scheduleData.value);
      }
    },
    fail: () => {
      uni.showToast({ title: 'Failed to load schedule', icon: 'none' });
    }
  });
};

const loadLogs = () => {
  uni.request({
    url: 'http://localhost:8080/logs',
    method: 'GET',
    success: (res) => {
      if (res.statusCode === 200) {
        logs.value = res.data;
      }
    },
    fail: () => {
      uni.showToast({ title: 'Failed to load logs', icon: 'none' });
    }
  });
};

const getCourseName = (row, col) => {
  const match = scheduleData.value.find(c => c.slot === row && c.day === col && currentWeek.value >= c.weekFrom && currentWeek.value <= c.weekTo);
  return match ? match.name : '[Empty]';
};

onMounted(() => {
  const savedList = uni.getStorageSync('classList') || [];
  classList.value = savedList;
  if (savedList.length > 0) {
    currentClass.value = savedList[0];
    loadSchedule();
    loadLogs();
  }
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
  width: 200px;
  border-right: 2px dashed #000;
  padding: 20px;
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
  padding: 10px 0;
  display: flex;
  justify-content: center;
  align-items: center;
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

.slot-btn {
  width: 100%;
  height: 100%;
  padding: 6px;
  background-color: #ffffff;
  border: 1px solid #ccc;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 0.95em;
  transition: background-color 0.2s;
}
.slot-btn:hover {
  background-color: #f0f0f0;
}
</style>
.schedule-table {
  width: 95%;
  max-width: 900px;
  margin: 0 auto;
  border-collapse: collapse;
  font-size: 0.95em;
  table-layout: fixed;
  text-align: center;
}

.schedule-table th,
.schedule-table td {
  border: 1px solid #000;
  padding: 12px;
  font-size: 1em;
}
.schedule-table thead {
  background-color: #f5f5f5;
}

.schedule-table th {
  background-color: #eaeaea;
  font-weight: bold;
  font-size: 1em;
}

.schedule-table td {
  height: 40px;
  vertical-align: middle;
}

.schedule-table td:first-child {
  font-weight: bold;
  background-color: #f8f8f8;
}