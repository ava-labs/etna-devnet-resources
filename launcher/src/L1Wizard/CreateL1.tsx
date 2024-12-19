import { useState, useEffect } from 'react';
import { useWizardStore } from "./store";
import NextPrev from "./ui/NextPrev";
import ReactConfetti from 'react-confetti';

type StepStatus = 'waiting' | 'loading' | 'success' | 'error';

export default function CreateL1() {
    const {
        genesisString,
        nodePopJsons,
        nodesCount,
        l1Name,
        chainId,
        subnetId,
        conversionId,
        setChainId,
        setSubnetId,
        setConversionId
    } = useWizardStore();

    const [isCreating, setIsCreating] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [status, setStatus] = useState<StepStatus>('waiting');
    const [showConfetti, setShowConfetti] = useState(false);
    const [windowSize, setWindowSize] = useState({
        width: window.innerWidth,
        height: window.innerHeight,
    });

    useEffect(() => {
        const handleResize = () => {
            setWindowSize({
                width: window.innerWidth,
                height: window.innerHeight,
            });
        };

        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    const handleCreate = async () => {
        setError(null);
        setIsCreating(true);
        setStatus('loading');

        try {
            const response = await fetch('/api/create', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    genesisString,
                    nodes: nodePopJsons.slice(0, nodesCount).map(json => JSON.parse(json).result),
                    l1Name,
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

            const data = await response.json();
            setSubnetId(data.subnetID);
            setChainId(data.chainID);
            setConversionId(data.conversionID);
            setStatus('success');
            setShowConfetti(true);
            setTimeout(() => setShowConfetti(false), 5000);
        } catch (error) {
            console.error('Creation error:', error);
            setStatus('error');
            setError(error instanceof Error ? error.message : 'Failed to create L1');
        }
    };

    const renderStepIcon = (status: StepStatus) => {
        switch (status) {
            case 'loading':
                return (
                    <svg className="animate-spin h-5 w-5 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                );
            case 'success':
                return (
                    <svg className="w-5 h-5 text-green-500" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 12">
                        <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M1 5.917 5.724 10.5 15 1.5" />
                    </svg>
                );
            case 'error':
                return (
                    <svg className="w-5 h-5 text-red-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                );
            case 'waiting':
                return (
                    <div className="w-5 h-5">
                        <svg className="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                            <path stroke="currentColor" strokeLinecap="round" strokeWidth="2" d="M6 12h.01m6 0h.01m5.99 0h.01" />
                        </svg>
                    </div>
                );
        }
    };

    const isCreationComplete = status === 'success';

    return (
        <div className="">
            {showConfetti && (
                <ReactConfetti
                    width={windowSize.width}
                    height={windowSize.height}
                    recycle={false}
                    numberOfPieces={800}
                    gravity={0.2}
                />
            )}
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
                <div className="flex flex-col mb-4">
                    <div className="flex items-center gap-3">
                        {renderStepIcon(status)}
                        <span className="text-gray-700">Create a Subnet</span>
                    </div>
                    {subnetId && status === 'success' && (
                        <p className="ml-8 mt-1 text-gray-600">
                            <code><a href={`https://subnets-test.avax.network/p-chain/tx/${subnetId}`} target="_blank" rel="noopener noreferrer" className="text-blue-500 underline hover:text-blue-700">{subnetId}</a></code>
                        </p>
                    )}
                </div>

                <div className="flex flex-col mb-4">
                    <div className="flex items-center gap-3">
                        {renderStepIcon(status)}
                        <span className="text-gray-700">Create a Chain</span>
                    </div>
                    {chainId && status === 'success' && (
                        <p className="ml-8 mt-1 text-gray-600">
                            <code><a href={`https://subnets-test.avax.network/p-chain/tx/${chainId}`} target="_blank" rel="noopener noreferrer" className="text-blue-500 underline hover:text-blue-700">{chainId}</a></code>
                        </p>
                    )}
                </div>

                <div className="flex flex-col mb-4">
                    <div className="flex items-center gap-3">
                        {renderStepIcon(status)}
                        <span className="text-gray-700">Convert the chain to an L1</span>
                    </div>
                    {conversionId && status === 'success' && (
                        <p className="ml-8 mt-1 text-gray-600">
                            <code><a href={`https://subnets-test.avax.network/p-chain/tx/${conversionId}`} target="_blank" rel="noopener noreferrer" className="text-blue-500 underline hover:text-blue-700">{conversionId}</a></code>
                        </p>
                    )}
                </div>
            </div>

            {!isCreationComplete && (
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
            )}

            <NextPrev
                nextDisabled={!isCreationComplete}
                currentStepName="create-l1"
            />
        </div>
    );
}
