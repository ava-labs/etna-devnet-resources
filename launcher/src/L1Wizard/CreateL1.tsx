import { useState } from 'react';
import { useWizardStore } from "./store";
import NextPrev from "./ui/NextPrev";

type StepStatus = 'waiting' | 'loading' | 'success' | 'error';

interface Step {
    label: string;
    status: StepStatus;
}

export default function CreateL1() {
    const { genesisString, nodePopJsons, nodesCount } = useWizardStore();
    const [isCreating, setIsCreating] = useState(false);
    const [steps, setSteps] = useState<Step[]>([
        { label: 'Create Subnet', status: 'waiting' },
        { label: 'Create Chain', status: 'waiting' },
        { label: 'Convert to L1', status: 'waiting' },
    ]);
    const [error, setError] = useState<string | null>(null);

    const handleCreate = async () => {
        setError(null);
        setIsCreating(true);
        setSteps(steps.map(step => ({ ...step, status: 'loading' as StepStatus })));

        try {
            const response = await fetch('/api/create', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    genesisString,
                    nodes: nodePopJsons.slice(0, nodesCount).map(json => JSON.parse(json).result),
                }),
            });

            if (!response.ok) {
                let errorMessage;
                try {
                    const errorData = await response.text();
                    errorMessage = errorData || response.statusText;
                } catch {
                    errorMessage = response.statusText;
                }
                throw new Error(errorMessage);
            }

            setSteps(steps.map(step => ({ ...step, status: 'success' as StepStatus })));
        } catch (error) {
            console.error('Creation error:', error);
            setSteps(steps.map(step => ({ ...step, status: 'error' as StepStatus })));
            setError(error instanceof Error ? error.message : 'Failed to create L1');
        }
    };

    return (
        <div className="">
            <h1 className="text-2xl font-medium mb-6">Create an L1</h1>

            {error && (
                <div className="mb-6 p-4 rounded-md bg-red-50 border border-red-200">
                    <div className="flex items-center gap-2">
                        <svg className="w-5 h-5 text-red-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span className="text-red-700 font-medium">Error</span>
                    </div>
                    <p className="mt-2 text-red-600">{error}</p>
                </div>
            )}
            <div className="mb-4">
                {steps.map((step, index) => (
                    <div key={index} className="flex items-center gap-3 mb-4">
                        {step.status === 'loading' && (
                            <svg className="animate-spin h-5 w-5 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                        )}
                        {step.status === 'success' && (
                            <svg className="w-5 h-5 text-green-500" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 12">
                                <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M1 5.917 5.724 10.5 15 1.5" />
                            </svg>
                        )}
                        {step.status === 'error' && (
                            <svg className="w-5 h-5 text-red-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        )}
                        {step.status === 'waiting' && (
                            <div className="w-5 h-5">
                                <svg className="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                                    <path stroke="currentColor" strokeLinecap="round" strokeWidth="2" d="M6 12h.01m6 0h.01m5.99 0h.01" />
                                </svg>
                            </div>
                        )}
                        <span className="text-gray-700">{step.label}</span>
                    </div>
                ))}
            </div>

            <div className="mb-8">
                <button
                    onClick={handleCreate}
                    disabled={isCreating}
                    className={`px-6 py-2 rounded-md ${isCreating
                        ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                        : 'bg-blue-500 text-white hover:bg-blue-600'
                        }`}
                >
                    {isCreating ? 'Creating...' : 'Create L1'}
                </button>
            </div>

            <NextPrev nextDisabled={false} currentStepName="create-l1" />
        </div>
    );
}
