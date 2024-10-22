<template>
  <span>Logout...</span>
</template>

<script>
import http from '../../http'

export default {
  name: 'Logout',
  mounted() {
    if (localStorage.getItem('go_token')) {
      this.logout()
    }
    else {
      this.login()
    }
  },
  methods: {
    logout() {
      http.get('/logout').finally(() => {
        localStorage.removeItem('go_token')
        this.login()
      })
    },
    login() {
      this.$router.push('login')
      this.$root.user = null
    }
  }
}
</script>