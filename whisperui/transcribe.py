import argparse
import whisper
import sys
import ffmpeg 
import numpy as np
import torch 
import whisper

def transcribe_audio(model_size: str, file_path: str) -> str:
    """
    Transcribes audio using the Whisper AI model.

    Args:
        model_size (str): The size of the Whisper model (e.g., "base", "small", "medium", "large").
        file_path (str): The path to the audio file to be transcribed.

    Returns:
        str: The transcribed text.
    """
    try:
        # Load the specified Whisper model
        model = whisper.load_model(model_size)
        
        # Transcribe the audio file
        result = model.transcribe(file_path)
        
        # Return the transcribed text
        return result['text']
    
    except Exception as e:
        return f"An error occurred: {e}"



if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Stream transcription from an audio file.")
    parser.add_argument("model_size", type=str, help="Whisper model size (tiny, base, small, medium, large)")
    parser.add_argument("filepath", type=str, help="Path to the audio file")

    args = parser.parse_args()

    sys.stdout.write(transcribe_audio(args.model_size, args.filepath))
    sys.exit(0)