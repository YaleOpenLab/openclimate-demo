pragma solidity ^0.5.10;

import "./ndcs.sol";
import "./reporting.sol";
// import "../math/SafeMath.sol";

contract ParisAgreementHighLevel{
    // using SafeMath for uint256;
    
    // At the high level we can set goals for every country and receive accumulated reports to keep track/actions

    /**
    Each country will be assigned its public ETH address and a private key kept secret on the responsible state agency
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
    function set_global_stocktake (address[] memory votingParties, int CO2, int CH4, int N2O, int AltEnergy, uint timeGoal) public pure {
        global_stocktake(votingPatries, CO2, CH4, N2O, AltEnergy, timeGoal);
    }
    
    
    /********************************************************************
    Action part
    *********************************************************************/

    // Import NDCs contract
    DetermedContributions contract_ndcs;
    function ImportNDCS(address _t) public {
        contract_ndcs = DetermedContributions(_t);
    }
    // Import Reporting contract
    Reporting report_ndc;
    function ReportNDCS(address _t) public {
        report_ndc = Reporting(_t);
    }
    // Calculate GHG reduction/surplus compared to NDC 
    function calculateReductions(address countryAddr) public view returns (int co2_reduction) {
        require(contract_ndcs.isCounrty(countryAddr), "Country doesnt have an NDC");
        (int ndcCO2, uint timeTarget) = contract_ndcs.getNdcCO2(countryAddr);
        (int reportCO2, uint timeStamp) = report_ndc.getLastCO2(countryAddr);
        require(timeTarget>=timeStamp, "Timestamp cant be higher than timeTarget");
        return (ndcCO2-reportCO2);
    }
    
    // //Make an action - issue/burn tokens depending on reporting
    // function ndcAction(address countryAddr) private (int ussuedTokens) {
    //    
    //    logic for issuing/burning tokens depending on how the country achieves its NDC target
    //    Token will probably have some value or benefits for a counrty (still working on it)
    //    
    //    return;
    // }
    
    /********************************************************************
     End of action part
    *********************************************************************/

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
}
