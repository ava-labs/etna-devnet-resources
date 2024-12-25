import { useWizardStore, resetStore } from "./store";
import { stepList } from "./store";

export default function Steps() {
    const { currentStep } = useWizardStore();
    const stepKeys = Object.keys(stepList) as (keyof typeof stepList)[];

    return (
        <>
            <ol className="relative text-gray-500 border-s border-gray-200">
                {stepKeys.map((stepKey, index) => {
                    const step = stepList[stepKey];
                    const isActive = stepKey === currentStep;
                    const isPast = stepKeys.indexOf(currentStep) > index;

                    return (
                        <li key={stepKey} className="mb-10 ms-6 last:mb-0">
                            <span className={`absolute flex items-center justify-center w-8 h-8 rounded-full -start-4 ring-2 ring-white 
                                ${isPast ? 'bg-green-200' :
                                    isActive ? 'bg-gray-100 ring-blue-500' :
                                        'bg-gray-100'}`}>
                                {isPast ? (
                                    <svg className="w-3.5 h-3.5 text-green-500" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 12">
                                        <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M1 5.917 5.724 10.5 15 1.5" />
                                    </svg>
                                ) : (
                                    <div className={`${isActive ? 'text-blue-500 outline outline-2 outline-blue-500 rounded-full p-1' : ''}`}>
                                        {step.icon}
                                    </div>
                                )}
                            </span>
                            <h3 className={`font-medium leading-tight ${isActive ? 'text-bold text-black' : ''}`}>
                                {step.title}
                            </h3>
                            <p className="text-sm">{step.description}</p>
                        </li>
                    );
                })}
            </ol>
            <div className="mt-8 -ml-4 w-full">
                <button
                    onClick={() => resetStore()}
                    className="flex items-center justify-center w-full px-4 py-2 text-sm font-medium text-red-600 bg-red-50 rounded-lg hover:bg-red-100"
                >
                    <svg className="w-4 h-4 mr-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
                    </svg>
                    Start Over
                </button>
            </div>
        </>
    );
}
