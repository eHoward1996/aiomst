<template>
  <div>
    <v-container fluid>
      Hello Darkness My Old Friend  <br>
    </v-container>
  </div>
</template>

<script>
export default {
  name: 'Index',
  mounted: function() {
    let stream = new EventSource("http://127.0.0.1:8090/stream")
    stream.addEventListener("msg", (e) => {
      console.log(e.data, Date().toLocaleLowerCase());
    });
    stream.addEventListener("close", () => {
      stream.close();
    })
  },
  data: function() {
    return {
      // active_tab: 0,
      // stream: new EventSource("http://127.0.0.1:8090/stream"),
    }
  },
  destroyed: function() {
    console.log('closed')
  }
}
</script>