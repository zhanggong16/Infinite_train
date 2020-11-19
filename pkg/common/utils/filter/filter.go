package filter

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
	"strconv"
)

//FiltersPara ...
type FiltersPara struct {
	FilterName      string
	FiltersOperator string   `json:"filters_operator"`
	FiltersValues   []string `json:"filters_values"`
}

//Filter ...
type Filter struct {
	IsPaging  bool          `json:"is_paging"`
	PageIndex int           `json:"page_index"`
	PageSize  int           `json:"page_size"`
	Filters   []FiltersPara `json:"filters"`
}

type FilterGroup struct {
	Filters []FiltersPost `json:"filters"`
}

type FiltersPost struct {
	Name     string   `json:"name"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

type SortPost struct {
	Name      string `json:"name"`
	Direction string `json:"direction"`
}

type FilterPostBody struct {
	FilterGroups []FilterGroup `json:"filter_groups"`
	Sort         SortPost      `json:"sort"`
	PageIndex    int           `json:"page_index"`
	PageSize     int           `json:"page_size"`
}

func (f *FilterPostBody) ConvertFromFilter(in Filter, c echo.Context) error {
	f.PageIndex = in.PageIndex
	f.PageSize = in.PageSize
	fgs := make([]FiltersPost, 0)
	for _, v := range in.Filters {
		fgs = append(fgs, FiltersPost{Name: v.FilterName, Operator: v.FiltersOperator, Values: v.FiltersValues})
	}
	f.FilterGroups = append(f.FilterGroups, FilterGroup{Filters: fgs})
	// sorts.1.name=lock_time&sorts.1.direction=asc， 因为日志服务只支持单列的排序，所以这里只处理第一个sort
	f.Sort.Name = c.QueryParam("sorts.1.name")
	f.Sort.Direction = c.QueryParam("sorts.1.direction")
	return nil
}

//GetFiltersFromURL ...
func (f *Filter) GetFiltersFromURL(c echo.Context) error {

	//**********************deal page**************************
	PageIndex := c.QueryParam("pageindex")
	PageSize := c.QueryParam("pagesize")
	f.IsPaging = true
	if PageIndex == "-1" {
		f.IsPaging = false

	} else if PageSize == "" && PageIndex == "" {
		f.IsPaging = false
	} else if PageSize != "" && PageIndex != "" {
		var err error
		var pageIndexInt int
		var pageSizeInt int
		validate := validator.New()
		pageIndexInt, err = strconv.Atoi(PageIndex)
		if err != nil {
			return err
		}
		if pageIndexInt <= 0 || pageIndexInt > 1000 {
			return errors.New("pageIndexInt <= 0 || pageIndexInt > 1000")
		}
		err = validate.Var(PageIndex, "number")
		if err != nil {
			return err
		}
		pageSizeInt, err = strconv.Atoi(PageSize)

		if err != nil {
			return err
		}
		if pageSizeInt <= 0 || pageSizeInt > 1000 {
			return errors.New("pageSizeInt <= 0 || pageSizeInt > 1000")
		}
		err = validate.Var(PageSize, "number")
		if err != nil {
			return err
		}
	} else {
		return errors.New("both of pageSize and pageNumber need or none")
	}

	f.PageIndex, _ = strconv.Atoi(PageIndex)
	f.PageSize, _ = strconv.Atoi(PageSize)

	//*********************deal filters***********************
	if c.QueryParam("filters.1.name") == "" {
		if c.QueryParam("filters.1.operator") != "" || c.QueryParam("filters.1.values.1") != "" {
			return fmt.Errorf("no name in filter 1, check the parameters carefully")
		}
		return nil
	}

	f.Filters = make([]FiltersPara, 0)
	var FiltersValues string
	var parameter FiltersPara
	var name, operator, values string

	for i := 1; ; i++ {
		//get name
		name = "filters." + strconv.Itoa(i) + ".name"
		parameter.FilterName = c.QueryParam(name)
		if parameter.FilterName == "" {
			operator = "filters." + strconv.Itoa(i) + ".operator"
			opTmp := c.QueryParam(operator)
			if opTmp != "" {
				return fmt.Errorf("no name in filter %d, check the parameters carefully", i)
			}
			values = "filters." + strconv.Itoa(i) + ".values.1"
			vaTmp := c.QueryParam(values)
			if vaTmp != "" {
				return fmt.Errorf("no name in filter %d, check the parameters carefully", i)
			}
			break
		}

		//get values
		for j := 1; ; j++ {
			values = "filters." + strconv.Itoa(i) + ".values." + strconv.Itoa(j)
			FiltersValues = c.QueryParam(values)
			if FiltersValues != "" {
				if IfExistStr(parameter.FiltersValues, FiltersValues) {
					return fmt.Errorf("value for Name:%s is repeated, check the parameters carefully", name)
				}
				parameter.FiltersValues = append(parameter.FiltersValues, FiltersValues)
			} else {
				break
			}
		}
		if len(parameter.FiltersValues) == 0 {
			return fmt.Errorf("name:%s, No values, check the parameters carefully", parameter.FilterName)
		}

		//get operator
		operator = "filters." + strconv.Itoa(i) + ".operator"
		parameter.FiltersOperator = c.QueryParam(operator)
		if parameter.FiltersOperator == "" {
			parameter.FiltersOperator = "eq"
			//return fmt.Errorf("Name:%s, No Operators", parameter.FilterName)
		}
		if parameter.FiltersOperator == "eq" && len(parameter.FiltersValues) != 1 {
			return fmt.Errorf("operator=[eq](if no operators, filter will use eq as default operator), but there are more than one values, check the parameters carefully")
		}
		f.Filters = append(f.Filters, parameter)
		parameter.FiltersValues = append([]string{})
	}
	return nil
}

// filters.1.name
// filters.1.operator
// filters.1.values.1
// filters.1.values.2
// filters.1.values.3
// .
// .
// .

//IfExistStr return if target exist in src
func IfExistStr(src []string, target string) bool {
	for _, s := range src {
		if s == target {
			return true
		}
	}
	return false
}
