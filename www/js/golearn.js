var synth = window.speechSynthesis;

function get_voices(lang) {
	voices = synth.getVoices();
	for (let l=0; l<voices.length; l++) {
		if (voices[l].lang == lang)
			return voices[l];
	}
	alert("Voices not found")
	return voices[0];
}

function say_word(word) {
	synth.cancel();
	if (synth.speaking) {
		console.error('speechSynthesis.speaking');
		return;
	}
	var ut = new SpeechSynthesisUtterance(word);
	ut.voice = get_voices("fr-FR");
	ut.rate = 0.7;
	synth.speak(ut);
}
