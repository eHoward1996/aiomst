<template>
  <v-row no-gutters align="center" justify="center">
    <v-col cols="1" align="center">
      <v-btn
          v-bind:disabled="!playback"
          icon
          @click="shuffle">
        <v-icon size="20">mdi-shuffle-variant</v-icon>
      </v-btn>
    </v-col>
    <v-col cols="1" align="center">
      <v-btn
          v-bind:disabled="!playback"
          icon
          @click="playPrev">
        <v-icon size="35">mdi-skip-previous</v-icon>
      </v-btn>
    </v-col>
    <v-col cols="1" align="center">
      <v-btn 
          v-if="!isPlaying"
          v-bind:disabled="!playback"
          icon
          @click="play">
        <v-icon size="50">mdi-play</v-icon>
      </v-btn>
      <v-btn 
          v-else
          icon
          @click="pause">
        <v-icon size="50">mdi-pause</v-icon>
      </v-btn>
    </v-col>
    <v-col cols="1" align="center">
      <v-btn
          v-bind:disabled="!playback"
          icon
          @click="playNext">
        <v-icon size="35">mdi-skip-next</v-icon>
      </v-btn>
    </v-col>
    <v-col cols="1" align="center">
      <v-btn
          v-bind:disabled="!playback"
          icon
          @click="loop">
        <v-badge
            color=""
            :icon="replayIcon"
            overlap
            bottom>
          <v-icon size="20">mdi-sync</v-icon>
        </v-badge>
      </v-btn>
    </v-col>
  </v-row>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
  name: 'controls',
  data: function() {
    return {
      progress: 0,
      sliderValue: 0,
      progressInterval: null,
    }
  },
  computed: {
    ...mapGetters({
      playback:     'currentSong',
      currentAlbum: 'currentAlbum',
      playlist:     'getPlaylist',
    }),
    isPlaying: function() {
      if (!this.playback) {
        return false
      }
      return this.playback.howl.playing(this.playback.howlId)
    },
    songDuration: function() {
      if (!this.playback) {
        return 0
      }
      return this.playback.song.length
    },
    replayIcon: function() {
      switch (this.$store.state.replayState) {
        case 0:
          return ''
        case 1:
          return  'mdi-numeric-1-circle'
        default:
          return 'mdi-infinity'
      }
    }
  },
  methods: {
    play: function() {
      var howl = this.playback.howl
      howl.play(this.playback.howlId)
    },
    pause: function() {
      var howl = this.playback.howl
      howl.pause(this.playback.howlId)
    },
    playNext: function() {
      var list = this.playlist
      var indexCurr = list.indexOf(this.playback.song)

      if (indexCurr + 1 >= list.length) {
        this.$store.dispatch('streamAudio', list[0])
          .then(() => console.log('next stream request done'))
        return
      }
      this.$store.dispatch('streamAudio', list[indexCurr + 1])
          .then(() => console.log('next stream request done'))
    },
    playPrev: function() {
      var list = this.playlist
      var indexCurr = list.indexOf(this.playback.song)

      if (indexCurr - 1 < 0) {
        this.$store.dispatch('streamAudio', list[0])
          .then(() => console.log('prev stream request done'))
        return
      }
      this.$store.dispatch('streamAudio', list[indexCurr - 1])
          .then(() => console.log('prev stream request done'))
    },
    shuffle: function() {
      console.log('shuffle')
    },
    loop: function()  {
      this.$store.commit('updateReplayState')
    }
  },
}
</script>

<style scoped lang="sass">
</style>