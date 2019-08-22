pragma solidity ^0.5.10;

import "./ndcs.sol";

contract ParisAgreementHighLevel{
    // At the high level we can set goals for every country and receive accumulated reports to keep track/actions

    /**
    Each country will be assigned is public ETH address and a private key kept secret on the responsible state agency
    We can rethink this process later.
    */
    address[] public countries;

    /**
    Global goals set by Conference of the Parties (COP). Related Voting contract will be separately added later.
    */
    struct global_stocktake {
        address[] votingPatries;
        int CO2; // (required) in metric tonnes
        int CH4; // (required) in metric tonnes
        int N2O; // (required) in metric tonnes
        int AltEnergy; // alternative/renewable energy usage in MWh
        uint timeGoal;
    }

    /**
     Util function set global goals by COP
     */
    function set_global_stocktale (address[] memory votingPatries, int CO2, int CH4, int N2O, int AltEnergy, uint timeGoal) public pure {
        global_stocktake(votingPatries, CO2, CH4, N2O, AltEnergy, timeGoal);
    }

    /**
    Probably we would need a intermediate total GHG variables to track who we do on the global level. Updated by information incoming from databases every year or month.
    */
    struct global_reduction_track {
        int  earthCO2;
        int  earthCH4;
        int  earthN2O;
        int  earthAltEnergy;
        uint  lastReportTime;
    }

    /**
     Util function to check time left to a time goal set up by COP (indicatively every 5 years)
     */
    // function checkDeadline () public returns (uint left, bool overdue){
    //     uint timeLeft = global_stocktake.timeGoal - now;
    //     if (timeLeft >= 0){
    //         return (timeLeft, false);
    //     }
    //     return (timeLeft, true);
    // }

    /**
     Import NDCs contract
     */
    DetermedContributions contract_ndcs = new DetermedContributions();
    /**
     Calculate CO2 reduction/surplus compared to NDC
     */
    function calcCO2Reduction(address country, int currentCO2) public view returns (int differnce) {
        if(!contract_ndcs.isCounrty(msg.sender)) revert();
        int ndcCO2 = contract_ndcs.getCO2(country);
        int diff = ndcCO2-currentCO2;
        return diff;
    }

}
