import 'htmx.org';
import 'theme-change';


import { VidstackPlayer, VidstackPlayerLayout } from 'vidstack/global/player';
import { MediaRemoteControl, TextTrack } from 'vidstack';
window.VidstackPlayer = VidstackPlayer;
window.VidstackPlayerLayout = VidstackPlayerLayout;
window.MediaRemoteControl = MediaRemoteControl;
window.TextTrack = TextTrack;
import 'vidstack/player';
import 'vidstack/player/layouts/default';
import 'vidstack/player/ui';



console.log("Eos")