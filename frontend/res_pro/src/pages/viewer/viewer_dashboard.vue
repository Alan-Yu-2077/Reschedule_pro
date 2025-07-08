<template>
  <div class="viewer-container">
    <div class="sidebar">
      <div class="class-header">
        <h2>Class List</h2>
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
                <div class="slot-display">{{ getCourseName(i, j) }}</div>
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
import { onLoad } from '@dcloudio/uni-app';

const classList = ref([]);
const currentClass = ref('');
const logs = ref([]);
const scheduleData = ref([]);

// Load class list
const loadClasses = async () => {
  try {
    const response = await uni.request({
      url: 'http://localhost:8080/api/schedule/classes',
      method: 'GET'
    });

    if (response.statusCode === 200) {
      classList.value = response.data.classes.map(cls => cls.name);
      console.log('Classes loaded:', classList.value);
      
      // Select first class if available
      if (classList.value.length > 0 && !currentClass.value) {
        currentClass.value = classList.value[0];
        loadSchedule();
      }
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

const currentWeek = ref(1);
const prevWeek = () => {
  if (currentWeek.value > 1) currentWeek.value--;
  loadSchedule();
};
const nextWeek = () => {
  if (currentWeek.value < 20) currentWeek.value++;
  loadSchedule();
};

// Load schedule data
const loadSchedule = async () => {
  if (!currentClass.value) return;
  
  try {
    const response = await uni.request({
      url: `http://localhost:8080/api/schedule/class/${encodeURIComponent(currentClass.value)}/week/${currentWeek.value}`,
      method: 'GET'
    });

    if (response.statusCode === 200) {
      scheduleData.value = response.data.schedules || [];
      console.log('Schedule loaded for week', currentWeek.value, ':', scheduleData.value);
    }
  } catch (error) {
    console.error('Failed to load schedule:', error);
    uni.showToast({ title: 'Failed to load schedule', icon: 'none' });
  }
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

// Get course name by time slot position
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

// Refresh data when page loads
onLoad(() => {
  loadClasses();
});
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Patrick+Hand&display=swap');

.viewer-container {
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
  padding: 8px 12px;
  margin-bottom: 5px;
  border: 1px solid #ddd;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.sidebar li:hover {
  background-color: #f0f0f0;
}

.schedule-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-left: 20px;
}

.schedule-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding: 10px;
  border-bottom: 2px dashed #000;
}

.controls button {
  padding: 8px 16px;
  border: 2px solid #000;
  background: #fff;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 1.1em;
  transition: transform 0.1s;
}

.controls button:hover {
  transform: translateY(-2px);
}

.title h2 {
  margin: 0;
  font-size: 1.5em;
}

.schedule-body {
  flex: 1;
  overflow: auto;
}

.schedule-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  border: 2px solid #000;
}

.schedule-table th,
.schedule-table td {
  border: 1px solid #ddd;
  padding: 12px 8px;
  text-align: center;
}

.schedule-table th {
  background: #f5f5f5;
  font-weight: bold;
  border-bottom: 2px solid #000;
}

.slot-display {
  padding: 8px;
  background: #f9f9f9;
  border-radius: 4px;
  min-height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9em;
}

.log-panel {
  margin-top: 20px;
  border-top: 2px dashed #000;
  padding-top: 20px;
}

.log-panel h3 {
  margin-bottom: 10px;
}

.log-content {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #ddd;
  padding: 10px;
  background: #f9f9f9;
}

.log-content div {
  margin-bottom: 5px;
  padding: 5px;
  border-bottom: 1px solid #eee;
}
</style> 