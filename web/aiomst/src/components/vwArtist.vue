<template>
  <div>
    <AppBar></AppBar>
    <Content page="artists"></Content>
  </div>
</template>

<script>
import AppBar from '@/components/cmpntAppBar.vue';
import Content from '@/components/cmpntContent.vue'

export default {
  name: 'Artists',
  components: {AppBar, Content},
}
</script>