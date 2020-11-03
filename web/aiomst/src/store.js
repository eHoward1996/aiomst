import axios from 'axios';
import Vue from 'vue';
import Vuex from 'vuex';
import {Howl} from 'howler';

Vue.use(Vuex);

export const store = new Vuex.Store({
  state: {
    gResp: null,
    playback: null,
    playlist: {},
    replayState: 0,
    shuffleState: false,
    universalVol: .5,
  },
  getters: {
    currentArtist: state => {
      return state.gResp['artists'].length === 1 ? state.gResp : null
    },      
    getArtists: state => {
      return state.gResp['artists'];
    },

    currentAlbum: state => {
      return state.gResp['albums'].length === 1 ? state.gResp : null
    },      
    getAlbums: state => {
      return state.gResp['albums'];
    },

    currentSong: state => {
      return state.playback
    },
    getSongs: state => {
      return state.gResp['songs'];
    },
    getPlaylist: state => {
      return state.shuffleState ? state.playlist.shuffled : state.playlist.origin;
    },

    getUniversalVol: state => {
      return state.universalVol;
    }
  },
  mutations: {
    changeGoResp: (state, payload) => {
      Vue.set(state, 'gResp', payload.gResp)
    },
    setSongState: (state, payload) => {
      Vue.set(state, 'playback', payload.playback)
    },
    setPlaylist: (state, payload) => {
      function randomize(a) {
        for (let i = a.length - 1; i > 0; i--) {
          const j = Math.floor(Math.random() * (i + 1));
          [a[i], a[j]] = [a[j], a[i]];
        }
        return a;
      }

      Vue.set(state.playlist, 'origin', payload)
      if (state.shuffleState) {
        var list = state.playlist.origin.slice()
        var indexCurr = list.indexOf(state.playback.song)
        var songCurr = list.splice(indexCurr, 1)
        
        Vue.set(state.playlist, 'shuffled', songCurr.concat(randomize(list)))
      }
    },
    updateReplayState: (state) => {
      Vue.set(
        state,
        'replayState',
        state.replayState === 2 ? 0 : state.replayState + 1)
    },
    updateShuffleState: (state) => {
      function randomize(a) {
        for (let i = a.length - 1; i > 0; i--) {
          const j = Math.floor(Math.random() * (i + 1));
          [a[i], a[j]] = [a[j], a[i]];
        }
        return a;
      }

      Vue.set(state, 'shuffleState', !state.shuffleState)
      if (state.shuffleState) {
        var list = state.playlist.origin.slice()
        var indexCurr = list.indexOf(state.playback.song)
        var songCurr = list.splice(indexCurr, 1)
        
        Vue.set(state.playlist, 'shuffled', songCurr.concat(randomize(list)))
      }
    },

    setUniversalVol: (state, payload) => {
      Vue.set(state, 'universalVol', payload);
    }
  },
  actions: {
    makeApiRequest: (context, navInfo) => {
      return new Promise((resolve, reject) => {
        axios
          .get(navInfo.path, {
            baseURL: 'http://127.0.0.1:8090',
            params: navInfo.params,
          })
          .then((x) => {
            context.commit('changeGoResp', {
              gResp: x.data,
            });
            resolve();
          })
          .catch((x) => {
            console.log(x);
            reject();
          });
      });
    },
    streamAudio: (context, song) => {
      // Check if a request for the current playing song has been made
      // If so, ignore it
      var playing = context.getters.currentSong;
      if (playing && playing.song.id === song.id) {
        return
      }
      
      return new Promise((resolve, reject) => {
        axios
          .get('http://localhost:8090/stream?id=' + song.id, {
            headers: {
              'Accept': 'text/event-stream',
            },
            responseType: 'arraybuffer'
          })
          .then(response => {
            // If there is a song playing, stop it
            if (context.getters.currentSong) {
              var playback = context.getters.currentSong
              playback.howl.unload()
            }
                  
            var blob = new Blob(
              [response.data], 
              {type: response.headers['content-type']}
            );

            var url = URL.createObjectURL(blob);
            var howlInfo = new Howl({
              src: [url],
              html5: true,
              volume: context.state.universalVol,
              preload: false,
              onend: () => {
                var list = context.getters.getPlaylist
                var indexCurr = list.indexOf(context.getters.currentSong.song)
                var next = list[0]

                switch (context.state.replayState) {
                  case 0:
                    // If the current song is the last in the playlist
                    // end playback
                    if (indexCurr === list.length - 1) {
                      return
                    }

                    // Else, get playback for the next song
                    next = list[indexCurr + 1]
                    break
                  case 2:
                    // If there is a 'next' song, create playback for it
                    if (indexCurr + 1 < list.length) {
                      next = list[indexCurr + 1]
                      break
                    }
                    // If the current song is the last song in a playlist with
                    // length > 1 use 'default' next value (list[0])
                    else if (list.length !== 1) {
                      break;
                    }
                    // Fallthrough if there is only 1 song in the playlist.
                  case 1:
                    // Stop (seek(0)) and start the song
                    var howlObj = context.state.playback.howl;
                    howlObj.stop()
                    howlObj.play(context.state.playback.howlId)
                    return
                }
                context.dispatch('streamAudio', next)
              }
            });
      
            var howlId = howlInfo.play();            
            context.commit(
              'setSongState', {
                'playback': {
                  'howl': howlInfo,
                  'howlId': howlId,
                  'song': song,
                }
              }
            )
            resolve();
          })
          .catch(error => {
            console.error(error);
            reject();
          });
      }); 
    }
  }
});