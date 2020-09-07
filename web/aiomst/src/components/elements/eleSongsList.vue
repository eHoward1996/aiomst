<template v-if="hasObjects">
  <v-container fluid>
    <v-list-item-group color="primary">
      <v-list-item
          v-for="song in getSongs"
          :key="song.id"
          @dblclick="playSong(song)"
          :ripple="false"
          color="white">
        <v-hover>
          <v-list-item-avatar
              @click="playSong(song)"
              slot-scope="{ hover }">
            <span v-if="isPlayback(song)">
              <v-progress-circular
                  :indeterminate="true"
                  :rotate="0"
                  :size="25"
                  :width="3"
                  color="green darken-2">
              </v-progress-circular>
            </span>
            <span v-else-if="!hover">{{formatTrackNumber(song)}}</span>
            <span v-else>
              <v-icon x-large color="green darken-2">mdi-play</v-icon>
            </span>
          </v-list-item-avatar>
        </v-hover>
        <v-list-item-content>
          <v-list-item-title v-html="song.title"></v-list-item-title>
          <v-list-item-subtitle v-html="song.artist"></v-list-item-subtitle>
        </v-list-item-content>
        <v-list-item-avatar>
          {{formatLength(song)}}
        </v-list-item-avatar>
      </v-list-item>
    </v-list-item-group>
  </v-container>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  name: 'eleSongsList',
  computed: {
    ...mapGetters({
      getSongs:    'getSongs',
      currentSong: 'currentSong',
      getPlaylist: 'getPlaylist',
    }),
    hasObjects() {
      return this.getSongs.length > 0;
    },
    hasPlayback() {
      return this.currentSong;
    },
  },
  methods: {
    formatLength: function(song) {
      var t = parseInt(song.length);
      var minute = Math.floor(t / 60);

      const zeroPad = (num) => String(num).padStart(2, '0')
      var sec = zeroPad(t % 60);

      return minute + ":" + sec;
    },
    formatTrackNumber: function(song) {
      const zeroPad = (num) => String(num).padStart(2, '0')
      return zeroPad(song.track)
    },
    isPlayback(song) {
      return this.hasPlayback && 
        this.currentSong.song.id === song.id &&
        this.currentSong.howl.playing(this.currentSong.howlId);
    },
    playSong: function(song) {
      this.$store.dispatch('streamAudio', song)
        .then(() => {
          console.log('stream request finished')
          this.$store.commit('setPlaylist', this.getSongs)
        })
        .catch((err) => {
          console.log(err)
        });
    },
  },
}
</script>