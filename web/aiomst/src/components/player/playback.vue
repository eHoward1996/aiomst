<template>
  <v-row 
      no-gutters
      class="playback">
    <v-col cols="1" class="text-center">
      {{ sliderValue ? formatTime(sliderValue) : '' }}
    </v-col>
    <v-col cols="10">
      <v-slider 
          :value="sliderValue" 
          :max="songDuration || 1"
          @change="manualSliderChange"
          hide-details>
      </v-slider>
    </v-col>
    <v-col cols="1" class="text-center">
      {{ songDuration ? formatTime(songDuration) : '' }}
    </v-col>
  </v-row>
</template>

<script>
import {mapGetters} from 'vuex';

export default {
  name: 'playback',
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
  },
  methods: {
    formatTime: function(timeToFormat) {
      var t = parseInt(timeToFormat);
      var minute = Math.floor(t / 60);

      const zeroPad = (num) => String(num).padStart(2, '0')
      var sec = zeroPad(t % 60);

      return minute + ":" + sec;
    },
    updateProgress: function() {
      if (!this.playback) {
        return
      }

      var pbHowl = this.playback.howl
      var pbHowlId = this.playback.howlId
      var durInMs = this.songDuration * 1000
      clearInterval(this.progressInterval);

      // If play Promise is not supported and we call internal actions the
      // playLock remains to true.
      // More Info: https://github.com/goldfire/howler.js/issues/1231
      if (pbHowl._playLock) {
        // Convert seconds to milliseconds (seek * 1000)
        this.progress = pbHowl._sounds[0]._seek * 1000
      } else {
        // Convert seconds to milliseconds (seek * 1000)
        this.progress = pbHowl.seek(null, pbHowlId) * 1000
      }

      if (pbHowl.playing(pbHowlId)) {
        this.progressInterval = setInterval(() => {
          if (this.progress + 1000 <= durInMs) {
            this.progress = this.progress + 1000;
          }
        }, 1000);
      }
    },
    manualSliderChange: function(sValue) {
      var howl = this.playback.howl
      var id = this.playback.howlId
      howl.seek(sValue, id)
    },
  },
  watch: {
    isPlaying() {
      this.updateProgress();
    },
    progress: function() {
      this.sliderValue = this.progress / 1000
    }
  }
}
</script>

<style scoped lang="sass">
  ::v-deep .v-slider
    height: 25px !important
    
    .v-slider__track-container
      height: 5px !important

      .v-slider__track-fill
        background: #388E3C !important

    .v-slider__thumb-container
      background: #888 !important

      .v-slider__thumb
        height: 15px
        width: 15px
        color: #AAAAAA !important
        background: #AAA !important

  .playback
    height: 50% !important
    align-items: center !important
    align-content: center !important
    justify-content: center !important

    .col
      width: 100% !important
      align-items: center !important
</style>