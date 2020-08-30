<template>
  <v-app>
    <cmpntAppBar></cmpntAppBar>
    <v-content>
      <v-container fluid>
        <router-view></router-view>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
import cmpntAppBar from '@/components/layout/cmpntAppBar.vue';

export default {
  name: 'App',
  components: {cmpntAppBar},
}
</script>