<template>
  <div> 
    <AppBar></AppBar>
    <Content page="album"></Content>
  </div>
</template>

<script>
import AppBar from '@/components/cmpntAppBar.vue';
import Content from '@/components/cmpntContent.vue'

export default {
  name: 'Album',
  components: {AppBar, Content},
}
</script>