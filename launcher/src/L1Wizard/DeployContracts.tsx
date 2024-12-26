import { useState, useEffect } from 'react';
import { useWizardStore } from './store';
import PoAValidatorManager from "../../contract_compiler/compiled/PoAValidatorManager.json"
import ValidatorMessages from "../../contract_compiler/compiled/ValidatorMessages.json"
import ProxyAdmin from "../../contract_compiler/compiled/ProxyAdmin.json"
import TransparentUpgradeableProxy from "../../contract_compiler/compiled/TransparentUpgradeableProxy.json"
import { createPublicClient, createWalletClient, custom, http, defineChain, keccak256 } from 'viem';
import NextPrev from './ui/NextPrev';

const PROXY_ADMIN_ADDRESS = '0xC0fFEE1234567890aBCdeF1234567890abcDef34' as const;
const PROXY_ADDRESS = '0x0Feedc0de0000000000000000000000000000000' as const;

type Status = 'not_started' | 'in_progress' | 'error' | 'success' | 'loading';

interface DeploymentStatus {
    status: Status;
    error?: string;
    address?: `0x${string}`;
    txHash?: `0x${string}`;
}

interface ProxyStatus {
    status: Status;
    error?: string;
    txHash?: `0x${string}`;
    currentImplementation?: `0x${string}`;
}

// Contract Deployment Card Component
function ContractCard({
    title,
    status,
    onDeploy,
    disabled = false,
    showDeployButton = true
}: {
    title: string;
    status: DeploymentStatus;
    onDeploy: () => Promise<void>;
    disabled?: boolean;
    showDeployButton?: boolean;
}) {
    const handleCopy = (text: string) => navigator.clipboard.writeText(text).catch(console.error);

    const statusColors = {
        not_started: 'bg-gray-50 border-gray-200',
        in_progress: 'bg-blue-50 border-blue-200',
        error: 'bg-red-50 border-red-200',
        success: 'bg-green-50 border-green-200',
        loading: 'bg-gray-50 border-gray-200'
    };

    return (
        <div className={`p-4 rounded-lg border ${statusColors[status.status]} mb-4`}>
            <div className="flex items-center justify-between mb-2">
                <h3 className="font-medium">{title}</h3>
                <span className={status.status === 'error' ? 'text-red-600' : ''}>
                    {status.status === 'not_started' ? 'Not deployed' :
                        status.status === 'in_progress' ? 'Deploying...' :
                            status.status === 'error' ? 'Failed' : 'Deployed'}
                </span>
            </div>

            {status.error && <div className="text-sm text-red-600 mb-2">{status.error}</div>}

            {status.address && (
                <div className="mb-2">
                    <div className="text-sm text-gray-500 mb-1">Contract Address:</div>
                    <div className="flex items-center bg-white rounded p-2 border border-gray-100">
                        <code className="font-mono text-sm flex-1 break-all">{status.address}</code>
                        <button onClick={() => handleCopy(status.address!)}
                            className="ml-2 px-2 py-1 text-xs bg-gray-100 hover:bg-gray-200 rounded">
                            Copy
                        </button>
                    </div>
                </div>
            )}

            {status.txHash && (
                <div>
                    <div className="text-sm text-gray-500 mb-1">Transaction Hash:</div>
                    <div className="flex items-center bg-white rounded p-2 border border-gray-100">
                        <code className="font-mono text-sm flex-1 break-all">{status.txHash}</code>
                        <button onClick={() => handleCopy(status.txHash!)}
                            className="ml-2 px-2 py-1 text-xs bg-gray-100 hover:bg-gray-200 rounded">
                            Copy
                        </button>
                    </div>
                </div>
            )}

            {showDeployButton && status.status === 'not_started' && (
                <button
                    onClick={onDeploy}
                    disabled={disabled}
                    className={`mt-2 w-full p-2 rounded ${disabled
                        ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                        : 'bg-blue-500 text-white hover:bg-blue-600'
                        }`}
                >
                    Deploy
                </button>
            )}
        </div>
    );
}

// Proxy Implementation Card Component
function ProxyCard({
    status,
    onUpgrade,
    disabled = false,
    validatorManagerAddress
}: {
    status: ProxyStatus;
    onUpgrade: () => Promise<void>;
    disabled?: boolean;
    validatorManagerAddress?: `0x${string}`;
}) {
    const handleCopy = (text: string) => navigator.clipboard.writeText(text).catch(console.error);

    const statusColors = {
        not_started: 'bg-gray-50 border-gray-200',
        in_progress: 'bg-blue-50 border-blue-200',
        error: 'bg-red-50 border-red-200',
        success: 'bg-green-50 border-green-200',
        loading: 'bg-gray-50 border-gray-200'
    };

    const needsUpgrade = status.currentImplementation !== validatorManagerAddress;

    return (
        <div className={`p-4 rounded-lg border ${statusColors[status.status]} mb-4`}>
            <div className="flex items-center justify-between mb-2">
                <h3 className="font-medium">Proxy Implementation</h3>
                <span className={status.status === 'error' ? 'text-red-600' : ''}>
                    {status.status === 'loading' ? 'Loading...' :
                        status.currentImplementation ? 'Current Implementation:' : 'Not Set'}
                </span>
            </div>

            {status.error && <div className="text-sm text-red-600 mb-2">{status.error}</div>}

            {status.currentImplementation && (
                <div className="mb-2">
                    <div className="flex items-center bg-white rounded p-2 border border-gray-100">
                        <code className="font-mono text-sm flex-1 break-all">{status.currentImplementation}</code>
                        <button onClick={() => handleCopy(status.currentImplementation!)}
                            className="ml-2 px-2 py-1 text-xs bg-gray-100 hover:bg-gray-200 rounded">
                            Copy
                        </button>
                    </div>
                </div>
            )}

            {validatorManagerAddress && needsUpgrade && (
                <button
                    onClick={onUpgrade}
                    disabled={disabled || status.status === 'loading'}
                    className={`mt-2 w-full p-2 rounded ${disabled || status.status === 'loading'
                        ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                        : 'bg-blue-500 text-white hover:bg-blue-600'
                        }`}
                >
                    Point to {validatorManagerAddress.slice(0, 6)}...{validatorManagerAddress.slice(-4)}
                </button>
            )}
        </div>
    );
}

// Main Component
export default function DeployContracts() {
    const [validatorMessages, setValidatorMessages] = useState<DeploymentStatus>({ status: 'not_started' });
    const [validatorManager, setValidatorManager] = useState<DeploymentStatus>({ status: 'not_started' });
    const [proxyStatus, setProxyStatus] = useState<ProxyStatus>({ status: 'loading' });

    const { evmChainId, rpcAddress, rpcLocationType, rpcDomainType, chainId } = useWizardStore();

    const getRpcEndpoint = () => {
        if (rpcLocationType === 'local') return `http://localhost:8080/ext/bc/${chainId}/rpc`;
        if (rpcDomainType === 'no-domain') return `https://${rpcAddress}.nip.io/ext/bc/${chainId}/rpc`;
        return `https://${rpcAddress}/ext/bc/${chainId}/rpc`;
    };

    const getChainConfig = () => defineChain({
        id: Number(evmChainId),
        name: chainId,
        nativeCurrency: {
            decimals: 18,
            name: 'Native Token',
            symbol: 'TOKEN',
        },
        rpcUrls: {
            default: { http: [getRpcEndpoint()] },
            public: { http: [getRpcEndpoint()] },
        },
    });

    const handleDeployValidatorMessages = async () => {
        setValidatorMessages({ status: 'in_progress' });
        try {
            if (!window.ethereum) {
                throw new Error('No wallet detected');
            }

            const chain = getChainConfig();
            const publicClient = createPublicClient({
                chain,
                transport: http(),
            });

            const walletClient = createWalletClient({
                chain,
                transport: custom(window.ethereum),
            });

            const [address] = await walletClient.requestAddresses();

            // Get bytecode and ensure it has 0x prefix
            const bytecode = ValidatorMessages.bytecode.object;
            const bytecodeWithPrefix = bytecode.startsWith('0x') ? bytecode : `0x${bytecode}`;

            // Deploy ValidatorMessages contract
            const hash = await walletClient.deployContract({
                chain,
                abi: ValidatorMessages.abi,
                bytecode: bytecodeWithPrefix as `0x${string}`,
                account: address,
            });

            // Wait for deployment and get contract address
            const receipt = await publicClient.waitForTransactionReceipt({ hash });

            if (!receipt.contractAddress) {
                throw new Error('Contract address not found in receipt');
            }

            setValidatorMessages({
                status: 'success',
                address: receipt.contractAddress,
                txHash: hash
            });
        } catch (err: any) {
            setValidatorMessages({
                status: 'error',
                error: err.message
            });
        }
    };

    const calculateLibraryHash = (libraryPath: string) => {
        // Calculate keccak256 of the fully qualified library name
        const hash = keccak256(
            new TextEncoder().encode(libraryPath)
        ).slice(2);
        // Take first 34 characters (17 bytes)
        return hash.slice(0, 34);
    };

    const handleDeployValidatorManager = async () => {
        if (!validatorMessages.address) {
            throw new Error('ValidatorMessages address not found');
        }

        setValidatorManager({ status: 'in_progress' });

        try {
            if (!window.ethereum) {
                throw new Error('No wallet detected');
            }

            const chain = getChainConfig();
            const publicClient = createPublicClient({
                chain,
                transport: http(),
            });

            const walletClient = createWalletClient({
                chain,
                transport: custom(window.ethereum),
            });

            const [address] = await walletClient.requestAddresses();

            // Get bytecode and ensure it has 0x prefix
            const bytecode = PoAValidatorManager.bytecode.object;
            const bytecodeWithPrefix = bytecode.startsWith('0x') ? bytecode : `0x${bytecode}`;

            // Get library path from linkReferences
            const libraryPath = `${Object.keys(PoAValidatorManager.bytecode.linkReferences)[0]}:${Object.keys(Object.values(PoAValidatorManager.bytecode.linkReferences)[0])[0]}`
            const libraryHash = calculateLibraryHash(libraryPath);
            const libraryPlaceholder = `__$${libraryHash}$__`;


            // Replace library placeholder with actual address using split/join
            const linkedBytecode = bytecodeWithPrefix
                .split(libraryPlaceholder)
                .join(validatorMessages.address.slice(2).padStart(40, '0'));

            if (linkedBytecode.includes("$__")) {
                throw new Error("Failed to replace library placeholder with actual address")
            }

            // Deploy ValidatorManager contract with linked bytecode
            const hash = await walletClient.deployContract({
                chain,
                abi: PoAValidatorManager.abi,
                bytecode: linkedBytecode as `0x${string}`,
                account: address,
            });

            // Wait for deployment and get contract address
            const receipt = await publicClient.waitForTransactionReceipt({ hash });

            if (!receipt.contractAddress) {
                throw new Error('Contract address not found in receipt');
            }

            setValidatorManager({
                status: 'success',
                address: receipt.contractAddress,
                txHash: hash
            });
        } catch (err: any) {
            setValidatorManager({
                status: 'error',
                error: err.message
            });
        }
    };

    useEffect(() => {
        checkProxyImplementation();
    }, []);

    const checkProxyImplementation = async () => {
        try {
            console.log("Checking proxy implementation");
            if (!window.ethereum) {
                throw new Error('No wallet detected');
            }

            const chain = getChainConfig();
            const publicClient = createPublicClient({
                chain,
                transport: http(),
            });

            // Create contract instance for the proxy
            const proxyContract = {
                address: PROXY_ADDRESS,
                abi: TransparentUpgradeableProxy.abi
            } as const;

            // Get the latest Upgraded event
            const events = await publicClient.getContractEvents({
                ...proxyContract,
                // eventName: 'Upgraded',
                // fromBlock: 'earliest',
                // toBlock: 'latest'
            });
            console.log("events", events);

            // if (events.length > 0) {
            //     // Get the most recent event
            //     const latestEvent = events[events.length - 1];
            //     const implementation = latestEvent.implementation;

            //     setProxyStatus({
            //         status: 'success',
            //         currentImplementation: implementation
            //     });
            // }
        } catch (err: any) {
            console.error('Failed to check proxy implementation:', err);
            setProxyStatus({
                status: 'error',
                error: err.message
            });
        }
    };

    const handleUpdateProxyAddress = async () => {
        if (!validatorManager.address) {
            throw new Error('ValidatorManager address not found');
        }

        setProxyStatus({ status: 'in_progress' });

        try {
            if (!window.ethereum) {
                throw new Error('No wallet detected');
            }

            const chain = getChainConfig();
            const publicClient = createPublicClient({
                chain,
                transport: http(),
            });

            const walletClient = createWalletClient({
                chain,
                transport: custom(window.ethereum),
            });

            const [address] = await walletClient.requestAddresses();

            // Create contract instance
            const proxyAdminContract = {
                address: PROXY_ADMIN_ADDRESS,
                abi: ProxyAdmin.abi,
            } as const;

            // Call upgradeAndCall function with empty data for a simple upgrade
            const hash = await walletClient.writeContract({
                ...proxyAdminContract,
                account: address,
                functionName: 'upgradeAndCall',
                args: [PROXY_ADDRESS, validatorManager.address, '0x'],
            });

            // Wait for transaction to complete
            const receipt = await publicClient.waitForTransactionReceipt({ hash });
            console.log("receipt", receipt);

            // Since we can't directly call getProxyImplementation, we'll update the state with the new implementation address
            setProxyStatus({
                status: 'success',
                txHash: hash,
                currentImplementation: validatorManager.address
            });
        } catch (err: any) {
            console.error('Proxy update failed:', err);
            setProxyStatus({
                status: 'error',
                error: err.message
            });
        }
    };

    return (
        <div className="max-w-3xl mx-auto">
            <h1 className="text-2xl font-medium mb-6">Deploy Contracts</h1>

            <ContractCard
                title="ValidatorMessages"
                status={validatorMessages}
                onDeploy={handleDeployValidatorMessages}
            />

            <ContractCard
                title="ValidatorManager"
                status={validatorManager}
                onDeploy={handleDeployValidatorManager}
                disabled={validatorMessages.status !== 'success'}
            />

            <ProxyCard
                status={proxyStatus}
                onUpgrade={handleUpdateProxyAddress}
                disabled={validatorManager.status !== 'success'}
                validatorManagerAddress={validatorManager.address}
            />

            <NextPrev
                nextDisabled={true}
                currentStepName="deploy-contracts"
            />
        </div>
    );
}
