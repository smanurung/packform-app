// initial component test
// TODO: add more tests.

import {shallowMount} from '@vue/test-utils'
import FilterBar from '@/components/FilterBar.vue'

describe('FilterBar', () => {
    it('correct default data', () => {
        expect(typeof FilterBar.data).toBe('function')
        const defaultData = FilterBar.data()
        expect(defaultData.filterText).toBe('')
        expect(defaultData.startDate).toBe('')
        expect(defaultData.endDate).toBe('')
    });
})