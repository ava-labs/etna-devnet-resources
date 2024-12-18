import { useWizardStore } from "./store";

export default function CreateL1() {
    const { advanceFrom } = useWizardStore();
    const nextDisabled = true

    return <>
        <h1 className="text-2xl font-medium mb-6">Create an L1</h1>
        <p className="mb-4">    WIP</p>
        <div className="flex justify-between">
            <button
                onClick={() => advanceFrom('create-l1', 'down')}
                className="px-4 py-2 rounded-md bg-gray-100 hover:bg-gray-200"
            >
                Previous
            </button>
            <button
                onClick={() => advanceFrom('create-l1')}
                disabled={nextDisabled}
                className={`px-4 py-2 rounded-md ${nextDisabled
                    ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                    : 'bg-blue-500 text-white hover:bg-blue-600'
                    }`}
            >
                Continue
            </button>
        </div>
    </>
}
