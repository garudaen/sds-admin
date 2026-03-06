/**
 * 应用状态：当前菜单、当前页面、当前选中的域名（用于解析记录页）
 */
export default {
  namespaced: true,
  state: {
    activeMenu: 'domain',
    currentPage: 'domain', // 'domain' | 'records'
    currentDomain: null,   // 进入解析记录页时选中的域名
  },
  mutations: {
    SET_ACTIVE_MENU(state, key) {
      state.activeMenu = key
    },
    SET_CURRENT_PAGE(state, page) {
      state.currentPage = page
    },
    SET_CURRENT_DOMAIN(state, domain) {
      state.currentDomain = domain
    },
    GO_TO_RECORDS(state, domain) {
      state.currentDomain = domain
      state.currentPage = 'records'
      state.activeMenu = 'records'
    },
    BACK_TO_DOMAINS(state) {
      state.currentDomain = null
      state.currentPage = 'domain'
      state.activeMenu = 'domain'
    },
  },
  actions: {
    setActiveMenu({ commit }, key) {
      commit('SET_ACTIVE_MENU', key)
      if (key === 'domain') {
        commit('SET_CURRENT_PAGE', 'domain')
        commit('SET_CURRENT_DOMAIN', null)
      }
    },
    goToRecords({ commit }, domain) {
      commit('GO_TO_RECORDS', domain)
    },
    backToDomains({ commit }) {
      commit('BACK_TO_DOMAINS')
    },
  },
  getters: {
    activeMenu: (s) => s.activeMenu,
    currentPage: (s) => s.currentPage,
    currentDomain: (s) => s.currentDomain,
  },
}
