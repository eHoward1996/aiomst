<template>
  <v-app>
    <AppBar></AppBar>
    <v-content>
      <v-container fluid>
        <router-view></router-view>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
import AppBar from '@/components/cmpntAppBar.vue';

export default {
  name: 'App',
  components: {AppBar},
  data() {
    return {}
  }
}
</script>