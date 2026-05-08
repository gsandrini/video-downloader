'use strict';
const TRANSLATIONS = {
    it: {
        title: 'Video Downloader',
        madeWith: 'Sviluppata con il supporto di',
        ytDlpMissing: 'YT-DLP non trovato. Installalo prima di procedere',
        labelUrl: 'URL video',
        placeholderUrl: 'https://youtu.be/...',
        labelStart: 'Inizio',
        labelEnd: 'Fine',
        labelOutput: 'File output',
        downloading: 'Download in corso...',
        btnDownload: 'Download',
    },
    en: {
        title: 'Video Downloader',
        madeWith: 'Developed with the support of',
        ytDlpMissing: 'YT-DLP not found. Install it before proceeding',
        labelUrl: 'Video URL',
        placeholderUrl: 'https://youtu.be/...',
        labelStart: 'Start',
        labelEnd: 'End',
        labelOutput: 'Output file',
        downloading: 'Download in progress...',
        btnDownload: 'Download',
    },
};

function VideoDownloader() {
    return {
        lang: navigator.language.startsWith('it') ? 'it' : 'en',
        get t() { return TRANSLATIONS[this.lang]; },
        async toggleLang() {
            this.lang = this.lang === 'it' ? 'en' : 'it';
            if (window.go?.main?.App) {
                await window.go.main.App.SetLanguage(this.lang);
            }
        },
        ytDlpAvailable: true,
        loading: false,
        result: null,
        form: {
            url: '',
            start: '',
            end: '',
            output: 'video.mp4',
        },

        async init() {
            if (window.go?.main?.App) {
                this.ytDlpAvailable = await window.go.main.App.CheckYtDlp()
                await window.go.main.App.SetLanguage(this.lang)
            }
        },

        async download() {
            this.result = null
            this.loading = true

            try {
                if (window.go?.main?.App) {
                    const res = await window.go.main.App.DownloadSegment(
                        this.form.url,
                        this.form.start,
                        this.form.end,
                        this.form.output || 'video.mp4'
                    )
                    this.result = res
                }
            } catch (e) {
                this.result = { success: false, message: String(e) }
            } finally {
                this.loading = false
            }
        },
    }
}