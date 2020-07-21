<template>
    <div class="filter-bar ui basic segment grid">
      <div class="ui form">
        <div class="inline field">
          <label>Search for:</label>
          <input type="text" v-model="filterText" class="three wide column" @keyup.enter="doFilter" placeholder="name, nickname, or email">
          <button class="ui primary button" @click="doFilter">Go</button>
        </div>
      </div>
      <div class="ui form">
        <div class="inline field">
          <label>Created date:</label>
          <input type="date" v-model="startDate" @change="doStartDate">
          <label>-</label>
          <input type="date" v-model="endDate" @change="doEndDate">
        </div>
      </div>
    </div>
</template>

<script>
  export default {
    data () {
      return {
        filterText: '',
        startDate: '',
        endDate: ''
      }
    },
    methods: {
      doFilter () {
        this.$events.fire('filter-set', this.filterText)
      },
      doStartDate () {
        this.$events.fire('start-date-set', this.startDate);
      },
      doEndDate() {
        this.$events.fire('end-date-set', this.endDate);
      }
    }
  }
</script>