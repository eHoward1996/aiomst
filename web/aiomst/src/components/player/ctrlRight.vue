<template>
  <v-container fluid v-if="playback">
    <v-row>
      <v-col cols="1">
        <v-icon>mdi-{{iconName}}</v-icon>
      </v-col>
      <v-col cols="6" class="sluder">
        <v-slider  
            :value="vol * 100"
            :max="100"
            @change="volChange"
            hide-details>
        </v-slider>
      </v-col>
      <v-col cols="1">
        <v-icon>mdi-volume-high</v-icon>
      </v-col>
      <v-col cols="4">
        Queue
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
  name: 'ctrlRight',
  data: function() {
    return {
      iconName: 'volume-medium',
      toggle: false,
    }
  },
  computed: {
    ...mapGetters({
      playback: 'currentSong',
      universalVol: 'getUniversalVol'
    }),
    vol: function() {
      return this.universalVol;
    },
  },
  methods: {
    toggleVol: function() {
      if (this.toggle || this.vol > 0) {
        this.toggle = !this.toggle
        // this.
      }

    },
    volChange: function(sValue) {
      var howl = this.playback.howl
      var id = this.playback.howlId
      howl.volume(sValue/100, id)
      this.$store.commit('setUniversalVol', (sValue/100))
    }
  },
  watch: {
    '$store.state.universalVol'() {
      var volume = this.vol * 100
      switch (typeof(volume) === 'number') {
        case volume === 0: 
          this.iconName = 'volume-off'
          break;
        case volume > 0 && volume <= 34:
          this.iconName = 'volume-low'
          break;
        case volume > 34 && volume <= 67:
          this.iconName = 'volume-medium'
          break
        case volume > 67 && volume <= 100:
          this.iconName = 'volume-high'
          break;
      }
    }
  }
}
</script>

<style lang="sass" scoped>
  ::v-deep .v-slider
    height: 2px !important
    
    .v-slider__track-container
      height: 5px !important

      .v-slider__track-fill
        background: #388E3C !important

    .v-slider__thumb-container
      background: #888 !important
      width: 70% !important

      .v-slider__thumb
        height: 10px
        width: 10px
        color: #AAAAAA !important
        background: #AAA !important
</style>