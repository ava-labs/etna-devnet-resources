import { useState } from 'react';
import { useWizardStore } from './store';
import NextPrev from "./ui/NextPrev";

type Status = 'not_started' | 'in_progress' | 'error' | 'success';

interface StepStatus {
    status: Status;
    error?: string;
    data?: any;
}

export default function CreateChain() {
    const { nodesCount } = useWizardStore();
    const [subnetStatus, setSubnetStatus] = useState<StepStatus>({ status: 'not_started' });
    const [createChainStatus, setCreateChainStatus] = useState<StepStatus>({ status: 'not_started' });

    const renderStepIcon = (status: Status) => {
        switch (status) {
            case 'in_progress':
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
            default:
                return (
                    <div className="w-5 h-5">
                        <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                            <path stroke="currentColor" strokeLinecap="round" strokeWidth="2" d="M6 12h.01m6 0h.01m5.99 0h.01" />
                        </svg>
                    </div>
                );
        }
    };

    const handleCreate = async () => {
        try {
            // Step 1: Simulate Subnet Creation
            setSubnetStatus({ status: 'in_progress' });
            try {
                await new Promise(resolve => setTimeout(resolve, 500));
                setSubnetStatus({
                    status: 'success',
                    data: 'subnet-' + Math.random().toString(36).substring(7)
                });
            } catch (error: any) {
                setSubnetStatus({
                    status: 'error',
                    error: error.message || 'Failed to create subnet'
                });
                return;
            }

            // Step 2: Simulate Chain Creation
            setCreateChainStatus({ status: 'in_progress' });
            try {
                await new Promise(resolve => setTimeout(resolve, 500));
                setCreateChainStatus({
                    status: 'success',
                    data: 'chain-' + Math.random().toString(36).substring(7)
                });
            } catch (error: any) {
                setCreateChainStatus({
                    status: 'error',
                    error: error.message || 'Failed to create chain'
                });
            }

        } catch (error: any) {
            console.error('Creation process failed:', error);
            if (subnetStatus.status === 'in_progress') {
                setSubnetStatus({
                    status: 'error',
                    error: 'Unexpected error during subnet creation'
                });
            } else if (createChainStatus.status === 'in_progress') {
                setCreateChainStatus({
                    status: 'error',
                    error: 'Unexpected error during chain creation'
                });
            }
        }
    };

    return (
        <div className="max-w-3xl mx-auto">
            <h1 className="text-2xl font-medium mb-6">Create a Chain</h1>

            {/* Display Selected Node Count */}
            <div className="mb-8">
                <h3 className="mb-4 font-medium">Selected Configuration</h3>
                <div className="p-4 bg-gray-50 rounded-lg">
                    <div className="flex items-center gap-2">
                        <span className="font-medium">{nodesCount} {nodesCount === 1 ? 'Node' : 'Nodes'}</span>
                        {nodesCount === 1 && <span className="bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded-full">Dev</span>}
                        {nodesCount === 3 && <span className="bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded-full">Testnet</span>}
                        {nodesCount === 5 && <span className="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full">Mainnet</span>}
                    </div>
                </div>
            </div>

            {/* Creation Steps */}
            <div className="mb-8 p-4 border rounded-lg">
                <h3 className="font-medium mb-4">Create Subnet & Chain</h3>
                <div className="mb-4">
                    <div className="flex flex-col mb-4">
                        <div className="flex items-center gap-3">
                            {renderStepIcon(subnetStatus.status)}
                            <span className="text-gray-700">Create a Subnet</span>
                        </div>
                        {subnetStatus.error && (
                            <p className="ml-8 mt-1 text-red-500 text-sm">
                                Error: {subnetStatus.error}
                            </p>
                        )}
                        {subnetStatus.data && (
                            <p className="ml-8 mt-1 text-gray-600">
                                <code><a href={`https://subnets-test.avax.network/p-chain/tx/${subnetStatus.data}`} target="_blank" rel="noopener noreferrer" className="text-blue-500 underline hover:text-blue-700">{subnetStatus.data}</a></code>
                            </p>
                        )}
                    </div>

                    <div className="flex flex-col mb-4">
                        <div className="flex items-center gap-3">
                            {renderStepIcon(createChainStatus.status)}
                            <span className="text-gray-700">Create a Chain</span>
                        </div>
                        {createChainStatus.error && (
                            <p className="ml-8 mt-1 text-red-500 text-sm">
                                Error: {createChainStatus.error}
                            </p>
                        )}
                        {createChainStatus.data && (
                            <p className="ml-8 mt-1 text-gray-600">
                                <code><a href={`https://subnets-test.avax.network/p-chain/tx/${createChainStatus.data}`} target="_blank" rel="noopener noreferrer" className="text-blue-500 underline hover:text-blue-700">{createChainStatus.data}</a></code>
                            </p>
                        )}
                    </div>
                </div>

                <button
                    onClick={handleCreate}
                    disabled={subnetStatus.status === 'in_progress' || createChainStatus.status === 'in_progress'}
                    className={`px-6 py-2 rounded-md ${subnetStatus.status === 'in_progress' || createChainStatus.status === 'in_progress'
                        ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                        : 'bg-blue-500 text-white hover:bg-blue-600'
                        }`}
                >
                    {subnetStatus.status === 'in_progress' || createChainStatus.status === 'in_progress'
                        ? 'Processing...'
                        : 'Start Process'}
                </button>
            </div>

            <NextPrev
                nextDisabled={createChainStatus.status !== 'success'}
                currentStepName="create-chain"
            />
        </div>
    );
}
