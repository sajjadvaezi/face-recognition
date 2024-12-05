import base64
import hashlib
import logging
import os
import time
from datetime import datetime

import cv2
import dlib
import numpy as np
import pickle
from flask import Flask, jsonify, request

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class FaceRecognitionApp:
    def __init__(self, storage_file="face_data.pkl", recognition_threshold=0.6):
        self.storage_file = storage_file
        self.recognition_threshold = recognition_threshold

        # Ensure necessary directories exist
        os.makedirs("pics", exist_ok=True)

        # Load face data
        self.face_data = self.load_face_data()

        # Initialize dlib's face detector (HOG-based) and facial landmark predictor
        try:
            self.detector = dlib.get_frontal_face_detector()
            self.sp = dlib.shape_predictor("shape_predictor_68_face_landmarks.dat")
            self.facerec = dlib.face_recognition_model_v1(
                "dlib_face_recognition_resnet_model_v1.dat"
            )
        except Exception as e:
            logging.error(f"Error initializing face recognition models: {e}")
            raise

    def load_face_data(self):
        """Load face data from a pickle file."""
        if os.path.exists(self.storage_file):
            try:
                with open(self.storage_file, "rb") as file:
                    return pickle.load(file)
            except (pickle.UnpicklingError, EOFError) as e:
                logging.error(f"Error loading face data: {e}")
                return {}
        return {}

    def save_face_data(self):
        """Save face data to a pickle file."""
        try:
            with open(self.storage_file, "wb") as file:
                pickle.dump(self.face_data, file)
        except Exception as e:
            logging.error(f"Error saving face data: {e}")

    def generate_unique_id(self, face_encoding):
        """Generate a unique ID by hashing the face encoding."""
        face_encoding_bytes = face_encoding.tobytes()
        return hashlib.sha256(face_encoding_bytes).hexdigest()

    def register_face(self, image_path):
        """Register a face and return a unique ID."""
        # Check if face is already recognized
        if (uid := self.recognize_face(image_path)) is not None:
            return uid

        # Load the image
        image = cv2.imread(image_path)
        if image is None:
            raise ValueError("Unable to read the image file")

        gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

        # Detect faces
        faces = self.detector(gray, 1)

        if len(faces) > 0:
            face = faces[0]
            shape = self.sp(image, face)
            face_encoding = np.array(self.facerec.compute_face_descriptor(image, shape))

            unique_id = self.generate_unique_id(face_encoding)
            self.face_data[unique_id] = face_encoding
            self.save_face_data()

            return unique_id
        else:
            raise ValueError("No face found in the image.")

    def recognize_face(self, image_path):
        """Recognize a face and return the unique ID or None."""
        # Load the image
        image = cv2.imread(image_path)
        if image is None:
            raise ValueError("Unable to read the image file")

        gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

        # Detect faces
        faces = self.detector(gray, 1)

        if len(faces) > 0:
            face = faces[0]
            shape = self.sp(image, face)
            face_encoding = np.array(self.facerec.compute_face_descriptor(image, shape))

            # Compare with stored encodings
            for unique_id, known_encoding in self.face_data.items():
                dist = np.linalg.norm(known_encoding - face_encoding)
                if dist < self.recognition_threshold:  # Configurable threshold
                    return unique_id
        return None

    def recognize_from_base64(self, base64_image):
        """Recognize face from base64 encoded image and save the image."""
        try:
            # Decode base64 image
            image_bytes = base64.b64decode(base64_image)
            nparr = np.frombuffer(image_bytes, np.uint8)
            image = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

            if image is None:
                raise ValueError("Invalid image data")
    
            # Save the uploaded image
            now = datetime.now()
            current_time = now.strftime("%Y%m%d_%H%M%S")
            saved_image_path = f"pics/recognize_upload_{current_time}.jpg"
            os.makedirs("pics", exist_ok=True)  # Ensure directory exists
            cv2.imwrite(saved_image_path, image)
            logging.info(f"Uploaded image saved to {saved_image_path}")

            # Convert to grayscale for face detection
            gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

            # Detect faces
            faces = self.detector(gray, 1)

            if len(faces) > 0:
                face = faces[0]
                shape = self.sp(image, face)
                face_encoding = np.array(self.facerec.compute_face_descriptor(image, shape))

                # Compare with stored encodings
                for unique_id, known_encoding in self.face_data.items():
                    dist = np.linalg.norm(known_encoding - face_encoding)
                    if dist < self.recognition_threshold:
                        return unique_id
            return None
        except Exception as e:
            logging.error(f"Error in recognize_from_base64: {e}")
            raise

    def register_from_base64(self, base64_image):
        """Register face from base64 encoded image and save the original image."""
        try:
            # Decode base64 image
            image_bytes = base64.b64decode(base64_image)
            nparr = np.frombuffer(image_bytes, np.uint8)
            image = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

            if image is None:
                raise ValueError("Invalid image data")

            # Save the original uploaded image
            now = datetime.now()
            current_time = now.strftime("%Y%m%d_%H%M%S")
            saved_image_path = f"pics/uploaded_{current_time}.jpg"
            os.makedirs("pics", exist_ok=True)  # Ensure directory exists
            cv2.imwrite(saved_image_path, image)
            logging.info(f"Uploaded image saved to {saved_image_path}")

            # Register the face from the saved image
            unique_id = self.register_face(saved_image_path)
            return unique_id
        except Exception as e:
            logging.error(f"Error in register_from_base64: {e}")
            raise



def take_photo(output_path='photo.jpg'):
    """Capture a photo from the default webcam."""
    cap = cv2.VideoCapture(0)

    try:
        if not cap.isOpened():
            logging.error("Could not open webcam")
            return False

        time.sleep(1)  # Camera warm-up
        ret, frame = cap.read()

        if ret:
            # Ensure pics directory exists
            os.makedirs("pics", exist_ok=True)
            cv2.imwrite(output_path, frame)
            logging.info(f"Photo saved as {output_path}")
            return True
        else:
            logging.error("Could not capture photo")
            return False
    except Exception as e:
        logging.error(f"Error in take_photo: {e}")
        return False
    finally:
        cap.release()


# Initialize the face recognition app
frApp = FaceRecognitionApp()

# Create Flask app
app = Flask(__name__)

@app.route('/health_check')
def health_check():
    """Simple health check endpoint."""
    return jsonify(status="everything is fine")

@app.route('/recognize', methods=["GET"])
def recognize():
    """Recognize a face from a captured photo."""
    now = datetime.now()
    current_time = now.strftime("%H%M%S")
    try:
        output_name = f"pics/{current_time}.jpg"
        if not take_photo(output_path=output_name):
            return jsonify(hash="", error="Failed to capture photo")

        res = frApp.recognize_face(output_name)
        if res is None:
            return jsonify(hash="", error="Couldn't recognize face")

        return jsonify(hash=res, error="")
    except Exception as e:
        logging.error(f"Error in recognize endpoint: {e}")
        return jsonify(hash="", error=str(e))

@app.route("/register", methods=["GET"])
def register_face():
    """Register a face from a captured photo."""
    now = datetime.now()
    current_time = now.strftime("%H%M%S")
    try:
        output_name = f"pics/{current_time}.jpg"
        if not take_photo(output_path=output_name):
            return jsonify(hash="", error="Failed to capture photo")

        hash_result = frApp.register_face(output_name)
        return jsonify(hash=hash_result, error="")
    except Exception as e:
        logging.error(f"Error in register endpoint: {e}")
        return jsonify(hash="", error=str(e))

@app.route('/recognize_upload', methods=['POST'])
def recognize_upload():
    """Recognize a face from an uploaded base64 image."""
    # Check if the request is JSON
    if not request.is_json:
        return jsonify(hash="", error="Request must be JSON"), 400

    # Get JSON data
    data = request.get_json()

    # Check if image is in the request
    if 'image' not in data:
        return jsonify(hash="", error="No image provided"), 400

    try:
        base64_image = data['image']
        # Validate base64 image
        if not base64_image:
            return jsonify(hash="", error="Empty image data"), 400

        # Recognize face
        res = frApp.recognize_from_base64(base64_image)

        if res is None:
            return jsonify(hash="", error="Couldn't recognize face"), 404

        return jsonify(hash=res, error="")
    except (ValueError, TypeError) as e:
        logging.error(f"Validation error in recognize_upload: {e}")
        return jsonify(hash="", error="Invalid image format"), 400
    except Exception as e:
        logging.error(f"Unexpected error in recognize_upload: {e}")
        return jsonify(hash="", error="Internal server error"), 500

@app.route('/register_upload', methods=['POST'])
def register_upload():
    """Register a face from an uploaded base64 image."""
    if 'image' not in request.json:
        return jsonify(hash="", error="No image provided"), 400

    try:
        base64_image = request.json['image']
        if not base64_image:
            return jsonify(hash="", error="Empty image data"), 400

        hash_result = frApp.register_from_base64(base64_image)
        return jsonify(hash=hash_result, error="")
    except ValueError as e:
        logging.error(f"Validation error in register_upload: {e}")
        return jsonify(hash="", error=str(e)), 400
    except Exception as e:
        logging.error(f"Unexpected error in register_upload: {e}")
        return jsonify(hash="", error="Internal server error"), 500

if __name__ == "__main__":
    app.run(debug=True)