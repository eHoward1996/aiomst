<template>
  <div> 
    <AppBar></AppBar>
    <Content page="albums"></Content>
  </div>
</template>

<script>
import AppBar from '@/components/cmpntAppBar.vue';
import Content from '@/components/cmpntContent.vue'

export default {
  name: 'Albums',
  components: {AppBar, Content},
}
</script>