import replicate
import dotenv
import sys
import json
dotenv.load_dotenv()


def main():
    if len(sys.argv) < 2:
        print("Usage: python stable_diff.py [options]")

    with open(sys.argv[1]) as f:
        data = json.load(f)
        predict(data)

def predict(data):
    model = replicate.models.get(f"{data['owner']}/{data['name']}")
    version = model.versions.get(data['version'])
    inputs = {
        'prompt': data['prompt'],
        'negative_prompt': data['neg_prompt'],
        'width': 512,
        'height': 512,
        'prompt_strength': 0.8,
        # Range: 1 to 4
        'num_outputs': data['num_outputs'],
        # Range: 1 to 500
        'num_inference_steps': 50,
        # Scale for classifier-free guidance
        # Range: 1 to 20
        'guidance_scale': 7.5,
        # Choose a scheduler.
        'scheduler': "DPMSolverMultistep",
    }
    output = version.predict(**inputs)
    print(json.dumps(output))


if __name__ == '__main__':
    main()
