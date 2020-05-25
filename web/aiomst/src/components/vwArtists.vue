<template>
  <!-- <v-app> -->
    <div>
    <AppBar></AppBar>
    <Content page="artists"></Content>
    </div>
  <!-- </v-app> -->
</template>

<script>
import AppBar from '@/components/cmpntAppBar.vue';
import Content from '@/components/cmpntContent.vue'

export default {
  name: 'Artists',
  components: {AppBar, Content},
}
</script>