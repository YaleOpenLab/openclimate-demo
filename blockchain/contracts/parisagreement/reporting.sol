pragma solidity ^0.5.10;

    /**
    Contract for to report intermediate accumulated GHG emissions from the countries Paris Agreement participants.
    Accumulated results from databases will commit automatically from Openclimate backend.
    */
contract Reporting {

    /**
    Struct to keep reporting entries
    */
    struct country_report {
        bytes32 country_name;
        // GHG mitigation goals
        int CO2; // (required) in metric tonnes
        int CH4; // (required) in metric tonnes
        int N2O; // (required) in metric tonnes
        int HFCs; // in metric tonnes
        int PFCs; // in metric tonnes
        int SF6; // in metric tonnes
        int NF3; // in metric tonnes
        int AltEnergy; // alternative/renewable energy usage in MWh
        // Finance
        int Mobilization; // (required) in billions USD.
        int ContribGreenFund; // (required) Contribution to Green Climate Fund in mibillionsllions USD
        int BilateralLoan; //  Bilateral loan to developing country Party billions USD
        // util variables
         timeStamp; // timestamp of the latest report
        uint index;
    }

    /**
    Mapping to store reports by address.
    */
    mapping(uint => country_report) private Reports;
    address[] private countryIndex;

    /**
    Util function to check if the address in the list
    */
    function isCounrty(address counrty)
    public
    view
    returns(bool isIndeed)
    {
        if(countryIndex.length == 0) return false;
        return (countryIndex[Reports[counrty].index] == counrty);
    }

    /**
       Function to insert a new report
    */
    function insert_report(
        bytes32 country_name,
        int CO2,
        int CH4,
        int N2O,
        int HFCs,
        int PFCs,
        int SF6,
        int NF3,
        int AltEnergy,
        int Mobilization,
        int ContribGreenFund,
        int BilateralLoan,
        uint timeStamp)
    public
    returns(uint index)
    {
        if(isCounrty(msg.sender)) revert();
        Reports[msg.sender].country_name      = country_name;
        Reports[msg.sender].CO2               = CO2;
        Reports[msg.sender].CH4               = CH4;
        Reports[msg.sender].N2O               = N2O;
        Reports[msg.sender].N2O               = HFCs;
        Reports[msg.sender].N2O               = PFCs;
        Reports[msg.sender].N2O               = SF6;
        Reports[msg.sender].N2O               = NF3;
        Reports[msg.sender].AltEnergy         = AltEnergy;
        Reports[msg.sender].Mobilization      = Mobilization;
        Reports[msg.sender].ContribGreenFund  = ContribGreenFund;
        Reports[msg.sender].BilateralLoan     = BilateralLoan;
        Reports[msg.sender].timeStamp         = timeStamp;
        Reports[msg.sender].index             = countryIndex.push(msg.sender)-1;
        return countryIndex.length-1;
    }
}
