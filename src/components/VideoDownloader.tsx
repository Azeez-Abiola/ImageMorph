import React, { useState } from 'react';
import { FaFacebook, FaYoutube, FaInstagram, FaTwitter, FaTiktok, FaLinkedin } from 'react-icons/fa';

const SUPPORTED_PLATFORMS = [
  { name: 'Facebook', icon: <FaFacebook className="text-blue-500" /> },
  { name: 'YouTube', icon: <FaYoutube className="text-red-500" /> },
  { name: 'Instagram', icon: <FaInstagram className="text-pink-500" /> },
  { name: 'Twitter', icon: <FaTwitter className="text-blue-500" /> },
  { name: 'TikTok', icon: <FaTiktok className="text-black" /> },
  { name: 'LinkedIn', icon: <FaLinkedin className="text-blue-500" /> },
];

const VideoDownloader = ({ isDarkMode }) => {
  const [url, setUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [videoData, setVideoData] = useState(null);
  const [selectedUrl, setSelectedUrl] = useState('');

 const handleFetch = async () => {
  setIsLoading(true);
  setError('');

  try {
    const response = await fetch('http://localhost:8080/api/metadata', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url: url.trim() }) 
    });

    if (!response.ok) throw new Error(await response.text());
    const data = await response.json();
    setVideoData(data);
    setSelectedUrl(data.medias?.[0]?.url || '');
  } catch (err) {
    setError(err.message);
  } finally {
    setIsLoading(false);
  }
};

const downloadVideo = () => {
  if (!selectedUrl) return;
  const safeTitle = videoData.Title?.replace(/[^\w]/g, '_') || 'video';
  window.open(
    `http://localhost:8080/api/download?url=${encodeURIComponent(selectedUrl)}&filename=${safeTitle}.mp4`
  )
}

 
  const inputStyles = isDarkMode
    ? 'bg-gray-900 text-white border-gray-700'
    : 'bg-white text-gray-900 border-gray-300';

  const buttonStyles = isLoading || !url
    ? 'opacity-50 cursor-not-allowed'
    : 'hover:bg-blue-600';

  const cardStyles = isDarkMode ? 'bg-gray-800' : 'bg-white';

  const formatDuration = (duration) => {
    const minutes = Math.floor(duration / 60);
    const seconds = duration % 60;
    return `${minutes}m ${seconds}s`;
  };

  return (
    <div onContextMenu={(e) => e.preventDefault()} className={`max-w-2xl mx-auto ${isDarkMode ? 'bg-gray-900 text-white' : 'bg-white text-gray-900'} select-none`}>
      <div className={`p-8 rounded-xl shadow-lg ${cardStyles}`}>
        <h2 className="text-2xl font-bold mb-6">Video Downloader</h2>
        <div className="flex items-center justify-center space-x-4 mb-6">
          {SUPPORTED_PLATFORMS.map((platform) => (
            <span key={platform.name} className="text-3xl">
              {platform.icon}
            </span>
          ))}
        </div>

        <div className="mb-8">
          <label htmlFor="videoUrl" className="font-medium block mb-2">
            Enter Video URL
          </label>
          <input
            type="url"
            id="videoUrl"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="Paste the video URL here..."
            className={`w-full px-4 py-3 rounded-lg border focus:ring-2 focus:ring-blue-500 ${inputStyles}`}
          />
          {error && <p className="text-red-500 text-sm mt-2">{error}</p>}
        </div>

        <button
          onClick={handleFetch}
          disabled={isLoading || !url}
          className={`w-full py-3 rounded-lg bg-blue-500 text-white font-medium ${buttonStyles}`}
        >
          {isLoading ? 'Fetching Video Info...' : 'Get Video'}
        </button>

        {videoData && (
          <div className="mt-6 p-4 rounded-lg border">
            <h3 className="text-lg font-bold mb-4">Video Details</h3>
            {videoData.thumbnail && (
              <img src={videoData.thumbnail} alt="Video Thumbnail" className="w-full rounded-md mb-4" />
            )}
            <p className={`${isDarkMode ? 'text-white' : 'text-gray-700'} mb-2`}>
              <strong>Title:</strong> {videoData.title}
            </p>
            <p className={`${isDarkMode ? 'text-white' : 'text-gray-700'} mb-2`}>
              <strong>Social Media Source:</strong> {videoData.source}
            </p>
            <p className={`${isDarkMode ? 'text-white' : 'text-gray-700'} mb-2`}>
              <strong>Duration:</strong> {formatDuration(videoData.duration)}
            </p>

            {videoData.medias && videoData.medias.length > 0 && (
              <div className="mt-4">
                <label htmlFor="qualitySelect" className="block mb-2">
                  Select Quality
                </label>
                <select
                  id="qualitySelect"
                  onChange={(e) => setSelectedUrl(e.target.value)}
                  className="w-full p-2 bg-black text-white rounded-md border"
                >
                  {videoData.medias.map((media, index) => (
                    <option key={index} value={media.url}>
                      Quality {media.quality}
                    </option>
                  ))}
                </select>
              </div>
            )}

            <button
              type="button"
              onClick={() => downloadVideo()}
              style={{
                pointerEvents: !selectedUrl ? 'none' : 'auto',
                opacity: !selectedUrl ? 0.5 : 1,
              }}
              className="cursor-pointer block w-full mt-4 bg-blue-500 text-center text-white p-3 rounded-md hover:bg-blue-600"
            >
              Download Video
            </button>
          </div>
        )}
      </div>
      <div className="mt-8 text-center text-sm text-gray-500">
        <p>Built with ❤️ by Abiola & Samuel</p>
      </div>
    </div>
  );
};

export default VideoDownloader;
