<template v-if="gResp['albums']">
  <div>
    <h1>Albums <v-icon color="primary" x-large>mdi-disc</v-icon></h1>
    <v-row 
      v-for="i in rowCount"
      :key="i">
      <v-col 
        v-for="album in gResp['albums'].slice((i - 1) * 3, i * 3)"
        :key="album.id">
        <v-card 
          max-width="300"
          max-height="450"
          v-on:click="albumAddr(album.id)">
          <v-img
            height="300"
            width="300"
            :src="artSrc(album)">
          </v-img>
          <v-card-title>{{album.title}}</v-card-title>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  name: 'rndAlbumsList',
  computed: {
    gResp() {
      return this.$store.state.gResp;
    },
    rowCount() {     
      return Math.ceil(this.gResp['albums'].length / 3);
    },
  },
  methods: {
    artSrc: function(album) {
      return "http://localhost:8090/art/" + album.artId;
    },
    albumAddr: function(album) {
      this.$router.push({
        path: 'album', 
        query: {id: album},
      })
    },
  },
}
</script>