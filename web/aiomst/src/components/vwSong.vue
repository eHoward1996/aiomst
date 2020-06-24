<template>
  <div>
    <AppBar></AppBar>
    <Content page="song"></Content>
  </div>
</template>

<script>
import AppBar from '@/components/cmpntAppBar.vue';
import Content from '@/components/cmpntContent.vue'

export default {
  name: 'Song',
  components: {AppBar, Content},
}
</script>