<template>
  <div class="ui container">
    <filter-bar></filter-bar>
    <vuetable ref="vuetable"
      api-url="http://localhost:8888/"
      :fields="fields"
      pagination-path=""
      @vuetable:pagination-data="onPaginationData"
      :append-params="moreParams"
    ></vuetable>
    <vuetable-pagination ref="pagination"
      @vuetable-pagination:change-page="onChangePage"></vuetable-pagination>
  </div>
</template>

<script>
import Vue from 'vue'
import Vuetable from 'vuetable-2/src/components/Vuetable'
import VuetablePagination from 'vuetable-2/src/components/VuetablePagination'
import FilterBar from './FilterBar'
import VueEvents from 'vue-events'

Vue.use(VueEvents)

export default {
  components: {
    Vuetable,
    VuetablePagination,
    FilterBar
  },
  mounted () {
    this.$events.$on('filter-set', eventData => this.onFilterSet(eventData));
    this.$events.$on('start-date-set', dt => this.onStartDateSet(dt));
    this.$events.$on('end-date-set', dt => this.onEndDateSet(dt));
  },
  methods: {
    onPaginationData (paginationData) {
      this.$refs.pagination.setPaginationData(paginationData)
    },
    onChangePage (page) {
      this.$refs.vuetable.changePage(page)
    },
    onFilterSet (filterText) {
      this.moreParams['filter'] = filterText;
      Vue.nextTick( () => this.$refs.vuetable.refresh());
    },
    onStartDateSet(dt) {
      this.moreParams['start_date'] = dt;
      Vue.nextTick(() => this.$refs.vuetable.refresh());
    },
    onEndDateSet(dt) {
      this.moreParams['end_date'] = dt;
      Vue.nextTick(() => this.$refs.vuetable.refresh());
    }
  },
  data () {
    return {
      fields: [
        {
          name: 'order_name',
          title: 'Order Name'
        },
        {
          name: 'customer_company',
          title: 'Customer Company'
        },
        {
          name: 'customer_name',
          title: 'Customer Name'
        },
        {
          name: 'order_date',
          title: 'Order date'
        },
        {
          name: 'delivered_amount',
          title: 'Delivered Amount'
        },
        {
          name: 'total_amount',
          title: 'Total Amount'
        }],
      moreParams: {
        per_page: 5
      }
    }
  }
}
</script>