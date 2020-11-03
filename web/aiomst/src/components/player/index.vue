<template>
  <v-container class="player" fluid pa-0>
    <v-row class="player_inner" no-gutters>
      <v-col cols="2" class="player_inner_left">
        <currentInfo></currentInfo>
      </v-col>
      <v-col cols="8" class="player_inner_center">
        <controls></controls>
        <playback></playback>
      </v-col>
      <v-col cols="2" class="player_inner_right">
        <ctrlRight></ctrlRight>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import currentInfo  from '@/components/player/currentInfo.vue';
import controls     from '@/components/player/controls.vue';
import playback     from '@/components/player/playback.vue';
import ctrlRight    from '@/components/player/ctrlRight.vue';

export default {
  name: 'elePlayerBar',
  components: {currentInfo, controls, playback, ctrlRight}
}
</script>

<style scoped lang="sass">
  .player
    position: fixed
    bottom: 0
    height: 90px
    z-index: 2
    background: #222

    &_inner
      height: 100%
      padding: .5%

      &_left
        height: 100%

      &_center
        height: 100%

      &_right
        height: 100%
        border: 1px solid cyan
</style>