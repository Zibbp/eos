// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/zibbp/eos/internal/video"

func VideoPlayer(video video.Video) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templ.JSONScript("video_path", video.VideoPath).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.JSONScript("thumbnail_path", video.ThumbnailPath).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.JSONScript("subtitle_path", video.SubtitlePath).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.JSONScript("ext_id", video.ExtID).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div><script type=\"module\">\n    // get variables set by server\n    const videoPath = JSON.parse(document.getElementById('video_path').textContent)\n    const thumbnailPath = JSON.parse(document.getElementById('thumbnail_path').textContent)\n    const subtitlePath = JSON.parse(document.getElementById('subtitle_path').textContent)\n    const extId = JSON.parse(document.getElementById('ext_id').textContent)\n    // create player\n    const player = await VidstackPlayer.create({\n      target: '#player',\n      src: [{ src: videoPath, type: 'video/webm' }],\n      poster: thumbnailPath,\n      layout: new VidstackPlayerLayout({\n        // thumbnails: 'https://files.vidstack.io/sprite-fight/thumbnails.vtt',\n      })\n    });\n\n    // create chapters track\n    player.textTracks.add({\n      type: 'json',\n      src: `/videos/${extId}/chapters`,\n      kind: 'chapters',\n      default: true\n    });\n    console.debug(player.textTracks)\n\n    // create subtitle tracks\n    if (subtitlePath && subtitlePath.length > 0) {\n    subtitlePath.forEach((subtitle) => {\n      console.log(subtitle)\n      const path = subtitle\n      const parts = subtitle.split('.')\n      const name = parts[parts.length - 2];\n      const type = parts[parts.length - 1]\n      const track = new TextTrack({\n        src: path,\n        kind: 'subtitles',\n        label: name,\n        language: name,\n        type: type\n      })\n\n      player.textTracks.add(track)\n    })\n    }\n\n    // set player volume from local storage\n    const localStorageVolume = localStorage.getItem(\"eos-player-volume\")\n    console.debug(`local volume: ${localStorageVolume}`);\n    if (localStorageVolume) {\n      player.volume = parseFloat(localStorageVolume);\n    }\n\n    // watch player volume saving to local storage\n    player.subscribe(({ volume }) => {\n      localStorage.setItem(\"eos-player-volume\", volume.toString())\n    })\n\n    // add class to player\n    const mediaPlayer = document.querySelector('media-player');\n    if (mediaPlayer) {\n      mediaPlayer.classList.add('eos-video-player');\n    }\n  </script><div id=\"player\"></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
