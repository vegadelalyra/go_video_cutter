const ffmpeg = require('fluent-ffmpeg');

// Define input and output files, and trimming start and duration times
const inputFile = 'C:/Users/rlyeh/Videos/2025-01-22 14-37-32.mkv';
const outputFile = '2025-01-22 14-37-32.mkv';
const startTime = '00:00:00'; // Start time in HH:MM:SS
const duration = 6_720; // Duration in seconds

// Trim the video
ffmpeg(inputFile)
  .setStartTime(startTime)
  .setDuration(duration)
  .output(outputFile)
  .on('end', () => {
    console.log('Video trimmed successfully!');
  })
  .on('error', err => {
    console.error('Error trimming video:', err.message);
  })
  .run();
