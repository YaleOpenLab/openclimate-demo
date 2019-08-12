pragma solidity ^0.5.10;

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
    According to Paris Agreement (PA) each country shall prepare, communicate and maintain successive nationally determined
    contributions (NDC) that it intends to achieve (GHG emissions). The struct allows to keep track of this for every country.

    */
    struct country_NDC {
        string country;
        // GHG mitigation goals
        int CO2; // (required) in metric tonnes
        int CH4; // (required) in metric tonnes
        int N2O; // (required) in metric tonnes
        int AltEnergy; // alternative/renewable energy usage in MWh

        // Adaptation goals


        // Capacity building


        // Finance
        int Mobilization; // (required) in billions USD
        int ContribGreenFund; // (required) Contribution to Green Climate Fund in mibillionsllions USD
        int BilateralLoan; //  Bilateral loan to developing country Party billions USD

        // Technology transfer


        // Transparency


        // util variables
        uint timeTarget; // set a time target
        uint index;
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
    Creating mapping to store countries goals. This allows to save/delete/update pledges for each country.
    */
    mapping(address => country_NDC) private NDCs;
    address[] private countryIndex;


    /**
    Util function to emit variable after inserting
    */
    event LogNewNDC   (address indexed country_address, uint index, string country_name, int CO2, int CH4, int N2O, int AltEnergy, uint timeDeadline);

    /**
    Util function toto check if the public key of the country in the list
    */
    function isCounrty(address counrty)
    public
    view
    returns(bool isIndeed)
    {
        if(countryIndex.length == 0) return false;
        return (countryIndex[Pledges[counrty].index] == counrty);
    }

    /**
       Function to insert a new NDC from a country address
    */
    function insertPledge(
        string memory country,
        int CO2,
        int CH4,
        int N2O,
        int AltEnergy,
        uint timeTarget)
    public
    returns(uint index)
    {
        if(isCounrty(msg.sender)) revert();
        NDCs[msg.sender].country           = country;
        NDCs[msg.sender].CO2               = CO2;
        NDCs[msg.sender].CH4               = CH4;
        NDCs[msg.sender].N2O               = N2O;
        NDCs[msg.sender].AltEnergy         = AltEnergy;
        NDCs[msg.sender].timeTarget      = timeTarget;
        NDCs[msg.sender].index             = countryIndex.push(msg.sender)-1;
        emit LogNewNDC(
            msg.sender,
            Pledges[msg.sender].index,
            country,
            CO2,
            CH4,
            N2O,
            AltEnergy,
                timeTarget);
        return countryIndex.length-1;
    }

    /**
     Util function to check time left to a time goal set up by COP (indicatively every 5 years)
     */
    function checkDeadline () public returns (uint time, bool overdue){
        uint timeLeft = global_stocktake.timeGoal - now;
        if (timeLeft >= 0){
            return (timeLeft, false);
        }
        return (timeLeft, true);
    }





}
