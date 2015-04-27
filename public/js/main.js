var CurveSignatures = new function () {
    var currentLayout;
    var elementHider;

    function Layout(name, maxRanks) {
        this.maxRanks = maxRanks;
        this.name = name;
    }

    function changedRankState(event) {
        if (currentLayout !== undefined) {
            if (getRanks(true).length > currentLayout.maxRanks) {
                event.target.checked = false;
                event.preventDefault();
                showAlert("Limit for amount of displayed ranks for selected layout is " + currentLayout.maxRanks);
            } else {
                clearAlert();
            }
        }
    }

    function clearAlert() {
        var element = document.getElementById("alert");
        if (element) {
            element.style.opacity = 0;
            elementHider = setTimeout(function () {
                element.parentNode.removeChild(element);
            }, 1000);
        }
    }

    function showAlert(message) {
        var element = document.getElementById("alert");

        if (elementHider) {
            clearTimeout(elementHider);
        }

        if (element) {
            element.innerHTML = message;
        } else {
            element = document.createElement("div");
            element.id = "alert";
            element.style.opacity = 0;
            element.innerHTML = message;
            var elements = document.getElementsByTagName("h3");
            elements[0].parentNode.appendChild(element);
            setTimeout(function () {
                element.style.opacity = 1.0;
            }, 100);
        }
        elementHider = setTimeout(function () {
            clearAlert();
        }, 3000);
    }

    function getRanks(selected) {
        var elements = Array.prototype.sort.call(Array.prototype.filter.call(document.getElementsByClassName("js-rank-selector"), function (value) {
            return selected && value.checked || !selected && !value.checked;
        }), function (a, b) {
            return b.getAttribute("data-rank") - a.getAttribute("data-rank");
        });
        return elements;
    }

    function changedLayout(event) {
        if (event !== null) {
            currentLayout = new Layout(event.target.getAttribute("name"), event.target.getAttribute("data-maxranks"));
        }
        var selectedRanks = getRanks(true);
        if (selectedRanks.length > currentLayout.maxRanks) {
            for (var i = 1; i <= (selectedRanks.length - currentLayout.maxRanks); i++) {
                selectedRanks[selectedRanks.length - i].checked = false;
            }
        } else {
            var ranks = getRanks(false);
            for (var i = 0; i < (currentLayout.maxRanks - selectedRanks.length); i++) {
                if (ranks[i] && ranks[i].getAttribute("data-rank") != 700) {
                    ranks[i].checked = true;
                }
            }
        }
    }

    this.init = function () {
        elementHider = setTimeout(function () {
            clearAlert();
        }, 3000);

        var elements = document.getElementsByClassName("js-layout-selector");
        if (elements.length > 0) {
            currentLayout = new Layout(elements[0].getAttribute("name"), elements[0].getAttribute("data-maxranks"));
            elements[0].checked = true;
            Array.prototype.forEach.call(elements, function (element) {
                element.addEventListener("change", changedLayout, false);
            });
            var elements = document.getElementsByClassName("js-rank-selector");
            Array.prototype.forEach.call(elements, function (element) {
                element.addEventListener("change", changedRankState, false);
            });
            changedLayout(null);
        }

        var elements = document.getElementsByClassName("js-select");
        Array.prototype.forEach.call(elements, function (element) {
            element.addEventListener("click", function (event) {
                if (window.getSelection().toString() === "") {
                    event.target.select();
                }
            });
        });



    };
};

CurveSignatures.init();