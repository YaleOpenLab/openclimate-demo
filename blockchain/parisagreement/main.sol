pragma solidity ^0.5.10;

contract ParisAgreementHighLevel{
    // At the high level we can set goals for every country and receive accumulated reports to keep track/actions

    /**
    Each country will be assigned is public ETH address and a private key kept secret on the responsible state agency
    */
    address[] public countries;

    /**
    According to Paris Agreement (PA) each country shall prepare, communicate and maintain successive nationally determined
    contributions that it intends to achieve (GHG emissions)
    */
    struct country_goal {
        string country;
        int CO2; // (required) in metric tonnes
        int CH4; // (required) in metric tonnes
        int N2O; // (required) in metric tonnes
        int AltEnergy; // alternative/renewable energy usage in MWh
        uint timeDeadline; // set a time target
        uint index;
    }

    /**
    Creating mapping to store countries goals
    */
    mapping(address => country_goal) private Pledges;
    address[] private countryIndex;

    event LogNewPledge   (address indexed country_address, uint index, string country_name, int CO2, int CH4, int N2O, int AltEnergy, uint timeDeadline);


    function isCounrty(address counrty)
    public
    view
    returns(bool isIndeed)
    {
        if(countryIndex.length == 0) return false;
        return (countryIndex[Pledges[counrty].index] == counrty);
    }

    /**
       Insert new pledge
    */
    function insertPledge(
        string memory country,
        int CO2,
        int CH4,
        int N2O,
        int AltEnergy,
        uint timeDeadline)
    public
    returns(uint index)
    {
        if(isCounrty(msg.sender)) revert();
        Pledges[msg.sender].country           = country;
        Pledges[msg.sender].CO2               = CO2;
        Pledges[msg.sender].CH4               = CH4;
        Pledges[msg.sender].N2O               = N2O;
        Pledges[msg.sender].AltEnergy         = AltEnergy;
        Pledges[msg.sender].timeDeadline      = timeDeadline;
        Pledges[msg.sender].index             = countryIndex.push(msg.sender)-1;
        emit LogNewPledge(
            msg.sender,
            Pledges[msg.sender].index,
            country,
            CO2,
            CH4,
            N2O,
            AltEnergy,
            timeDeadline);
        return countryIndex.length-1;
    }


    /**
    According to Paris Agreement (PA) each country shall communicate a nationally determined contribution every five years
    */


}
