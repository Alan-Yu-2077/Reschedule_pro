<template>
  <div class="main-container">
    <div class="title-area">
      <h2>Register New Class Schedule</h2>
      <div class="class-name-block">
        <label>Class Name:</label>
        <input v-model="className" placeholder="e.g. Grade 23 - Class 1" class="class-name-input" />
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
              <button class="slot-btn" @click="openDialog(i, j)">
                {{ schedule[i][j]?.name || '[Empty]' }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showDialog" class="dialog-overlay">
      <div class="dialog-box">
        <h3>Add Course</h3>
        <label>Course Name:</label>
        <input v-model="newCourse.name" placeholder="e.g. Math" />
        <label>Start Week:</label>
        <input type="number" v-model="newCourse.startWeek" />
        <label>End Week:</label>
        <input type="number" v-model="newCourse.endWeek" />
        <div class="dialog-actions">
          <button @click="confirmCourse">Confirm</button>
          <button @click="closeDialog">Cancel</button>
        </div>
      </div>
    </div>
    <div class="submit-area">
      <button class="submit-btn" @click="submitSchedule">Submit</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'

const schedule = ref(
  Array(5).fill(null).map(() => Array(7).fill(null))
)

const className = ref('')

onShow(() => {
  className.value = ''
  schedule.value = Array(5).fill(null).map(() => Array(7).fill(null))
})

const submitSchedule = () => {
  if (!className.value.trim()) {
    uni.showToast({ title: 'Class name is required', icon: 'none' });
    return;
  }

  const payload = {
    className: className.value,
    weeks: {
      1: schedule.value.map(row =>
        row.map(cell => (cell ? cell.name : ""))
      )
    }
  };

  uni.request({
    url: 'http://localhost:8080/schedule',
    method: 'POST',
    data: payload,
    success: (res) => {
      if (res.statusCode === 200) {
        const existing = uni.getStorageSync('classList') || []
        if (!existing.includes(className.value)) {
          existing.push(className.value)
        }
        uni.setStorageSync('classList', existing)
        uni.showToast({ title: 'Schedule submitted', icon: 'success' });
        uni.navigateTo({ url: '/pages/Main_page/main' });
      } else {
        uni.showToast({ title: 'Submit failed', icon: 'none' });
      }
    },
    fail: () => {
      uni.showToast({ title: 'Request failed', icon: 'none' });
    }
  });
}

const showDialog = ref(false)
const selectedRow = ref(0)
const selectedCol = ref(0)
const newCourse = ref({
  name: '',
  startWeek: 1,
  endWeek: 1
})

const openDialog = (row, col) => {
  selectedRow.value = row
  selectedCol.value = col
  newCourse.value = { name: '', startWeek: 1, endWeek: 1 }
  showDialog.value = true
}

const closeDialog = () => {
  showDialog.value = false
}

const confirmCourse = () => {
  schedule.value[selectedRow.value][selectedCol.value] = { ...newCourse.value }
  closeDialog()
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
  max-width: 1000px;
  margin: 0 auto;
}

.title-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 10px;
}

.schedule-body {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 20px;
  width: 100%;
}

.schedule-table {
  width: 95%;
  max-width: 900px;
  border-collapse: collapse;
  table-layout: fixed;
  text-align: center;
}

.schedule-table th,
.schedule-table td {
  border: 1px solid #000;
  padding: 12px;
}

.schedule-table th {
  background-color: #eaeaea;
}

.schedule-table td:first-child {
  font-weight: bold;
  background-color: #f8f8f8;
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

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.3);
  display: flex;
  justify-content: center;
  align-items: center;
}

.dialog-box {
  background: white;
  border: 2px dashed #000;
  padding: 20px;
  font-family: 'Patrick Hand', cursive;
  width: 300px;
}

.dialog-box h3 {
  margin-bottom: 10px;
}

.dialog-box label {
  display: block;
  margin-top: 10px;
  font-size: 0.9em;
}

.dialog-box input {
  width: 100%;
  padding: 6px;
  margin-top: 4px;
  box-sizing: border-box;
  font-family: 'Patrick Hand', cursive;
}

.dialog-actions {
  margin-top: 15px;
  display: flex;
  justify-content: space-between;
}
.dialog-actions button {
  padding: 6px 10px;
  font-family: 'Patrick Hand', cursive;
  cursor: pointer;
}

.class-name-input {
  margin-top: 4px;
  padding: 6px;
  width: 300px;
  font-family: 'Patrick Hand', cursive;
}

.submit-area {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  width: 100%;
  max-width: 900px;
}

.submit-btn {
  padding: 8px 16px;
  font-family: 'Patrick Hand', cursive;
  font-size: 1em;
  cursor: pointer;
  border: 1px solid #000;
  background: #fff;
  transition: background 0.2s;
}
.submit-btn:hover {
  background: #f5f5f5;
}
</style>
.class-name-block {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 8px;
}