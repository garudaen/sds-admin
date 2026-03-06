/**
 * 域名列表状态（可选：用于在添加/编辑后刷新列表）
 */
export default {
  namespaced: true,
  state: {
    list: [],
  },
  mutations: {
    SET_LIST(state, list) {
      state.list = list || []
    },
  },
  actions: {
    setList({ commit }, list) {
      commit('SET_LIST', list)
    },
  },
  getters: {
    list: (s) => s.list,
  },
}
