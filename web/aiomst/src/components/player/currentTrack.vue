<template>
  <v-container class="player" fluid pa-0>
    <v-row class="player_inner" no-gutters>

      <v-col 
          v-if="playback" 
          cols="2"
          order="1"
          class="player_inner_left">

        <v-row no-gutters style="height: 100%;">
          <v-col cols="5" class="curr-img">
            <v-row no-gutters style="height: 100%;" align-content="center">
              <v-col cols="8">
                <v-img :src="getArtSrc()" width="90px" height="70px"></v-img>            
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
      </v-col>
      <v-col 
          v-else 
          cols="2" 
          order="1"
          class="player_inner_left">
      </v-col>

      <v-col 
          cols="8"
          order="2"
          class="player_inner_center">
        <v-row no-gutters style="height: 50%;" class="controls">
          <v-col cols="3"></v-col>
          <v-col cols="6" align="center">
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
          <v-col cols="3"></v-col>
        </v-row>

        <v-row no-gutters style="height: 50%;" class="playback">
          <v-col cols="2" class="outer text-center">
            {{ sliderValue ? formatTime(sliderValue) : '' }}
          </v-col>
          
          <v-col cols="8" class="inner" align-self="center">
            <v-slider 
                :value="sliderValue" 
                :max="songDuration || 1"
                @change="manualSliderChange">
            </v-slider>
          </v-col>
          
          <v-col cols="2" class="outer text-center">
            {{ songDuration ? formatTime(songDuration) : '' }}
          </v-col>
        </v-row>
      </v-col>

      <v-col 
          cols="2"
          order="3"
          class="player_inner_right">
        <v-row no-gutters>right</v-row>
      </v-col>

    </v-row>
  </v-container>
</template>

<script>
import {mapGetters} from 'vuex';

export default {
  name: 'elePlayerBar',
  data: function() {
    return {
      progress: 0,
      sliderValue: 0,
      progressInterval: null,
      prevArtId: 0,
    }
  },
  computed: {
    ...mapGetters({
      playback: 'currentSong',
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
    play: function() {
      var howl = this.playback.howl
      howl.play(this.playback.howlId)
    },
    pause: function() {
      var howl = this.playback.howl
      howl.pause(this.playback.howlId)
    },
    formatTime: function(timeToFormat) {
      var t = parseInt(timeToFormat);
      var minute = Math.floor(t / 60);

      const zeroPad = (num) => String(num).padStart(2, '0')
      var sec = zeroPad(t % 60);

      return minute + ":" + sec;
    },
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
      console.log(sValue)
    }
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

  .player
    position: fixed
    bottom: 0
    height: 90px
    z-index: 2
    background: #333

    &_inner
      height: 100%
      padding: .5% 

      &_left
        height: 100%
        
        .curr-img
          .col
            min-width: 0
            height: 100%
            margin: auto

        .curr-info 
          .col
            min-width: 0
            overflow: hidden
            white-space: nowrap
            text-overflow: ellipsis
            font-size: 14px
            margin: auto
            margin-right: 15%

      &_center
        height: 100%

        .playback
          .inner
            height: 100%
            align-self: center !important

          .outer
            width: 100%
            align-self: center !important


      &_right
        border: 1px solid cyan
        height: 100%
</style>