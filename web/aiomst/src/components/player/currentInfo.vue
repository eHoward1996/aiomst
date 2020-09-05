<template>
  <v-row v-if="playback" no-gutters>
    <v-col cols="4" class="curr-img">
      <v-row no-gutters align-content="center">
        <v-col>
          <v-img :src="getArtSrc()" width="80px" height="70px"></v-img>            
        </v-col>
      </v-row>
    </v-col>
    
    <v-col 
        cols="7"
        class="curr-info"
        align-self="center">
      <v-row no-gutters>
        <v-col>{{ playback.song.title  }}</v-col>
      </v-row>
      <v-row no-gutters>
        <v-col>{{ playback.song.artist }}</v-col>
      </v-row>
      <v-row no-gutters>
        <v-col>{{ playback.song.album  }}</v-col>
      </v-row>
    </v-col>
  </v-row>
</template>
  
<script>
import { mapGetters } from 'vuex'
export default {
  name: 'currentInfo',
  data: function() {
    return {
      prevArtId: 0,
    }
  },
  computed: {
    ...mapGetters({
      playback:     'currentSong',
      currentAlbum: 'currentAlbum',
    })
  },
  methods: {
    getArtSrc: function() {
      if (this.currentAlbum) {
        var album = this.currentAlbum["albums"][0]
        var artId = album.artId
        if (album.id === this.playback.song.albumId) {
          this.prevArtId = artId
          return 'http://localhost:8090/art?id=' + artId
        }
      }
      return 'http://localhost:8090/art?id=' + this.prevArtId
    },
  }
}
</script>

<style scoped lang="sass">
  .curr-info
    .col
      min-width: 0
      overflow: hidden
      white-space: nowrap
      text-overflow: ellipsis
      font-size: 14px
</style>