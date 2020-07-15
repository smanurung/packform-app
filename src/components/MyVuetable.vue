<template>
  <div class="ui container">
    <filter-bar></filter-bar>
    <vuetable ref="vuetable"
      api-url="https://vuetable.ratiw.net/api/users"
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

Vue.component('filter-bar', FilterBar)
Vue.use(VueEvents)

export default {
  components: {
    Vuetable,
    VuetablePagination
  },
  mounted () {
    this.$events.$on('filter-set', eventData => this.onFilterSet(eventData))
    this.$events.$on('filter-reset', () => this.onFilterReset())
  },
  methods: {
    onPaginationData (paginationData) {
      this.$refs.pagination.setPaginationData(paginationData)
    },
    onChangePage (page) {
      this.$refs.vuetable.changePage(page)
    },
    onFilterSet (filterText) {
      this.moreParams = {
        'filter': filterText
      }
      Vue.nextTick( () => this.$refs.vuetable.refresh())
    },
    onFilterReset () {
      this.moreParams = {}
      Vue.nextTick( () => this.$refs.vuetable.refresh())
    }
  },
  data () {
    return {
      fields: ['name', 'email', 'birthdate',
        {
          name: 'address.line1',
          title: 'Address 1'
        },
        {
          name: 'address.line2',
          title: 'Address 2'
        },
        {
          name: 'address.zipcode',
          title: 'Zipcode'
        }],
      moreParams: {}
    }
  }
}
</script>