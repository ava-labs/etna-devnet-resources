import { useEffect, useState } from 'react';
import { createPublicClient, createWalletClient, custom, http, defineChain } from 'viem';
import { useWizardStore } from './store';
import ProxyAdmin from "../../contract_compiler/compiled/ProxyAdmin.json";
import TransparentUpgradeableProxy from "../../contract_compiler/compiled/TransparentUpgradeableProxy.json";
import PoAValidatorManager from "../../contract_compiler/compiled/PoAValidatorManager.json";

const PROXY_ADMIN_ADDRESS = '0xC0fFEE1234567890aBCdeF1234567890abcDef34' as const;
const PROXY_ADDRESS = '0x0Feedc0de0000000000000000000000000000000' as const;
const IMPLEMENTATION_SLOT = '0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc' as const;

export const ProxyDebug = () => {
    const [proxyLogs, setProxyLogs] = useState<any[]>([]);
    const [proxyAdminLogs, setProxyAdminLogs] = useState<any[]>([]);
    const [proxyOwner, setProxyOwner] = useState<string>('');
    const [error, setError] = useState<string>('');
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


    const handleUpgradeProxy = async () => {
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

            const newImplementation = "0x745591a9acd18d0eff15bb6cdfd8e05449b3f554"
            console.log('Upgrading proxy to implementation:', newImplementation);

            const hash = await walletClient.writeContract({
                address: PROXY_ADMIN_ADDRESS,
                abi: ProxyAdmin.abi,
                functionName: 'upgradeAndCall',
                args: [PROXY_ADDRESS, newImplementation, '0x'],
                account: address,
            });

            const receipt = await publicClient.waitForTransactionReceipt({ hash });
            console.log('Upgrade transaction receipt:', receipt);

            console.log('Upgrade transaction hash:', hash);
        } catch (err: any) {
            console.error('Failed to upgrade proxy:', err);
            setError(err.message);
        }
    };

    useEffect(() => {
        const fetchData = async () => {
            try {
                const chain = getChainConfig();
                const publicClient = createPublicClient({
                    chain,
                    transport: http(),
                });

                const implementationSlotData = await publicClient.getStorageAt({
                    address: PROXY_ADDRESS,
                    slot: IMPLEMENTATION_SLOT,
                });

                console.log('Current Implementation:', implementationSlotData);

                console.log("contractCode at PROXY_ADDRESS", await publicClient.getCode({
                    address: PROXY_ADDRESS,
                }));

                console.log("contractCode at PROXY_ADMIN_ADDRESS", await publicClient.getCode({
                    address: PROXY_ADMIN_ADDRESS,
                }));

                // Fetch UPGRADE_INTERFACE_VERSION
                const upgradeVersion = await publicClient.readContract({
                    address: PROXY_ADMIN_ADDRESS,
                    abi: ProxyAdmin.abi,
                    functionName: 'UPGRADE_INTERFACE_VERSION',
                });
                console.log("UPGRADE_INTERFACE_VERSION:", upgradeVersion);

                // Fetch owner
                const owner = await publicClient.readContract({
                    address: PROXY_ADMIN_ADDRESS,
                    abi: ProxyAdmin.abi,
                    functionName: 'owner',
                });
                setProxyOwner(String(owner));

                // Fetch Proxy logs
                const proxyLogsContract = {
                    address: PROXY_ADDRESS,
                    abi: TransparentUpgradeableProxy.abi,
                };
                const proxyEvents = await publicClient.getContractEvents({
                    ...proxyLogsContract,
                    fromBlock: 0n,
                    toBlock: 'latest',
                });

                // Fetch ProxyAdmin logs
                const proxyAdminContract = {
                    address: PROXY_ADMIN_ADDRESS,
                    abi: ProxyAdmin.abi,
                };
                const proxyAdminEvents = await publicClient.getContractEvents({
                    ...proxyAdminContract,
                    fromBlock: 0n,
                    toBlock: 'latest',
                });

                // Check PoAValidatorManager functions
                const [userAddress] = await window.ethereum.request({
                    method: 'eth_requestAccounts'
                });

                const proxyContract = {
                    address: PROXY_ADDRESS,
                    abi: PoAValidatorManager.abi,
                };

                console.log('Checking PoAValidatorManager functions through proxy:');

                const addressLength = await publicClient.readContract({
                    ...proxyContract,
                    functionName: 'ADDRESS_LENGTH',
                    account: userAddress as `0x${string}`,
                });
                console.log('ADDRESS_LENGTH:', addressLength);

                const blsKeyLength = await publicClient.readContract({
                    ...proxyContract,
                    functionName: 'BLS_PUBLIC_KEY_LENGTH',
                });
                console.log('BLS_PUBLIC_KEY_LENGTH:', blsKeyLength);

                const maxChurn = await publicClient.readContract({
                    ...proxyContract,
                    functionName: 'MAXIMUM_CHURN_PERCENTAGE_LIMIT',
                });
                console.log('MAXIMUM_CHURN_PERCENTAGE_LIMIT:', maxChurn);

                setProxyLogs(proxyEvents);
                setProxyAdminLogs(proxyAdminEvents);
            } catch (err: any) {
                setError(err.message);
                console.error('Error fetching data:', err);
            }
        };

        fetchData();
    }, []);

    return (
        <div className="max-w-3xl mx-auto p-4">
            <h1 className="text-2xl font-bold mb-4">Proxy Debug</h1>

            {error && (
                <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
                    Error: {error}
                </div>
            )}

            <div className="mb-6">
                <button
                    onClick={handleUpgradeProxy}
                    className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                >
                    Upgrade Proxy (Random Implementation)
                </button>
            </div>

            <div className="mb-6">
                <h2 className="text-xl font-semibold mb-2">Proxy Admin Owner</h2>
                <div className="bg-gray-100 p-4 rounded">
                    <code className="font-mono">{proxyOwner || 'Loading...'}</code>
                </div>
            </div>

            <div className="mb-6">
                <h2 className="text-xl font-semibold mb-2">Proxy Logs</h2>
                <pre className="bg-gray-100 p-4 rounded overflow-auto max-h-96">
                    {JSON.stringify(proxyLogs, null, 2)}
                </pre>
            </div>

            <div>
                <h2 className="text-xl font-semibold mb-2">Proxy Admin Logs</h2>
                <pre className="bg-gray-100 p-4 rounded overflow-auto max-h-96">
                    {JSON.stringify(proxyAdminLogs, null, 2)}
                </pre>
            </div>
        </div>
    );
};
