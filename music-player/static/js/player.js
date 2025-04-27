document.addEventListener('DOMContentLoaded', function() {
    const audioPlayer = document.getElementById('audio-player');
    const playBtn = document.getElementById('play-btn');
    const prevBtn = document.getElementById('prev-btn');
    const nextBtn = document.getElementById('next-btn');
    const modeSelect = document.getElementById('mode-select');
    const songList = document.getElementById('song-list');
    const currentSongTitle = document.getElementById('current-song');
    
    let currentSongIndex = 0;
    let isPlaying = false;
    let playerMode = 'sequence';
    
    // 获取当前播放器状态
    fetch('/api/toggle')
        .then(response => response.json())
        .then(data => {
            if (data.status === 'playing') {
                isPlaying = true;
                playBtn.textContent = '暂停';
            }
        });
    
    // 获取当前播放模式
    fetch('/api/mode?mode=' + modeSelect.value)
        .then(response => response.json())
        .then(data => {
            playerMode = data.mode;
        });
    
    // 播放/暂停按钮
    playBtn.addEventListener('click', function() {
        const action = isPlaying ? 'pause' : 'play';
        fetch('/api/toggle')
            .then(response => response.json())
            .then(data => {
                if (data.status === 'playing') {
                    audioPlayer.play();
                    playBtn.textContent = '暂停';
                    isPlaying = true;
                } else {
                    audioPlayer.pause();
                    playBtn.textContent = '播放';
                    isPlaying = false;
                }
            });
    });
    
    // 上一首按钮
    prevBtn.addEventListener('click', function() {
        fetch('/api/prev')
            .then(response => response.json())
            .then(data => {
                loadSong(data);
            });
    });
    
    // 下一首按钮
    nextBtn.addEventListener('click', function() {
        fetch('/api/next')
            .then(response => response.json())
            .then(data => {
                loadSong(data);
            });
    });
    
    // 模式选择
    modeSelect.addEventListener('change', function() {
        playerMode = this.value;
        fetch('/api/mode?mode=' + playerMode)
            .then(response => response.json())
            .then(data => {
                console.log('播放模式已切换为:', data.mode);
            });
    });
    
    // 加载歌曲
    function loadSong(songData) {
        // 在实际应用中，这里应该从API获取歌曲数据
        // 这里简化处理，直接使用模拟数据
        const songs = [
            {Title: "Song 1", Artist: "Artist A", Path: "/static/music/song1.mp3"},
            {Title: "Song 2", Artist: "Artist B", Path: "/static/music/song2.mp3"},
            {Title: "Song 3", Artist: "Artist C", Path: "/static/music/song3.mp3"},
            {Title: "Song 4", Artist: "Artist D", Path: "/static/music/song4.mp3"}
        ];
        
        // 根据模式获取下一首歌曲
        let nextIndex;
        switch(playerMode) {
            case 'random':
                nextIndex = Math.floor(Math.random() * songs.length);
                break;
            case 'repeat':
                nextIndex = currentSongIndex;
                break;
            case 'sequence':
            default:
                nextIndex = (currentSongIndex + 1) % songs.length;
                break;
        }
        
        // 更新当前索引
        currentSongIndex = nextIndex;
        
        const song = songs[currentSongIndex];
        
        // 更新UI
        audioPlayer.src = song.Path;
        currentSongTitle.textContent = `${song.Title} - ${song.Artist}`;
        
        // 高亮当前歌曲
        const songItems = songList.querySelectorAll('li');
        songItems.forEach((item, index) => {
            if (index === currentSongIndex) {
                item.classList.add('active');
            } else {
                item.classList.remove('active');
            }
        });
        
        // 播放歌曲
        if (isPlaying) {
            audioPlayer.play();
        }
    }
    
    // 点击歌曲列表项
    songList.addEventListener('click', function(e) {
        if (e.target.tagName === 'LI') {
            const songItem = e.target;
            const songPath = songItem.getAttribute('data-path');
            
            // 找到对应的歌曲索引
            const songs = [
                {Title: "Song 1", Artist: "Artist A", Path: "/static/music/song1.mp3"},
                {Title: "Song 2", Artist: "Artist B", Path: "/static/music/song2.mp3"},
                {Title: "Song 3", Artist: "Artist C", Path: "/static/music/song3.mp3"},
                {Title: "Song 4", Artist: "Artist D", Path: "/static/music/song4.mp3"}
            ];
            
            let index = -1;
            songs.forEach((song, i) => {
                if (song.Path === songPath) {
                    index = i;
                }
            });
            
            if (index !== -1) {
                currentSongIndex = index;
                
                // 高亮选中的歌曲
                const songItems = songList.querySelectorAll('li');
                songItems.forEach((item, i) => {
                    if (i === index) {
                        item.classList.add('active');
                    } else {
                        item.classList.remove('active');
                    }
                });
                
                // 加载并播放歌曲
                audioPlayer.src = songPath;
                currentSongTitle.textContent = `${songs[index].Title} - ${songs[index].Artist}`;
                
                if (isPlaying) {
                    audioPlayer.play();
                }
            }
        }
    });
    
    // 音频播放结束事件
    audioPlayer.addEventListener('ended', function() {
        fetch('/api/next')
            .then(response => response.json())
            .then(data => {
                loadSong(data);
            });
    });
});