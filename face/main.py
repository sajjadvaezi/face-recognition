import cv2
import dlib
import pickle
import hashlib
import numpy as np
import os


class FaceRecognitionApp:
    def __init__(self, storage_file="face_data.pkl"):
        self.storage_file = storage_file
        self.face_data = self.load_face_data()

        # Initialize dlib's face detector (HOG-based) and facial landmark predictor
        self.detector = dlib.get_frontal_face_detector()
        self.sp = dlib.shape_predictor("shape_predictor_68_face_landmarks.dat")
        self.facerec = dlib.face_recognition_model_v1(
            "dlib_face_recognition_resnet_model_v1.dat"
        )

    def load_face_data(self):
        """Load face data from a pickle file."""
        if os.path.exists(self.storage_file):
            with open(self.storage_file, "rb") as file:
                return pickle.load(file)
        return {}

    def save_face_data(self):
        """Save face data to a pickle file."""
        with open(self.storage_file, "wb") as file:
            pickle.dump(self.face_data, file)

    def generate_unique_id(self, face_encoding):
        """Generate a unique ID by hashing the face encoding."""
        face_encoding_bytes = face_encoding.tobytes()
        return hashlib.sha256(face_encoding_bytes).hexdigest()

    def register_face(self, image_path):
        """Register a face and return a unique ID."""
        if (uid:=self.recognize_face(image_path)) is not None:
            return uid
        # Load the image
        image = cv2.imread(image_path)
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
                if dist < 0.6:  # Threshold for matching
                    return unique_id
        return None


if __name__ == "__main__":
    app = FaceRecognitionApp()

    try:
        unique_id = app.register_face('c.jpg')  # Change to your image path
        print(f"Registered face with ID: {unique_id}")
    except ValueError as e:
        print(e)

    recognized_id = app.recognize_face("c.jpg")  # Change to your image path
    if recognized_id is not None:
        print(f"Recognized face with ID: {recognized_id}")
    else:
        print("Face not recognized.")


# e85c8420e15a03ff4cadc929ca344feb7b4b90ab98a2af1588bab1c39d9b5e5c
# e85c8420e15a03ff4cadc929ca344feb7b4b90ab98a2af1588bab1c39d9b5e5c
# e85c8420e15a03ff4cadc929ca344feb7b4b90ab98a2af1588bab1c39d9b5e5c
# e85c8420e15a03ff4cadc929ca344feb7b4b90ab98a2af1588bab1c39d9b5e5c